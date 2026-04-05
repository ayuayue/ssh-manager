package keys

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SSHKey struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Type       string `json:"type"`
	Size       int    `json:"size"`
	IsPrivate  bool   `json:"isPrivate"`
	PubKeyContent string `json:"pubKeyContent"`
	Fingerprint string `json:"fingerprint"`
}

type KeyGenRequest struct {
	Type     string `json:"type"`
	Bits     int    `json:"bits"`
	Email    string `json:"email"`
	Passphrase string `json:"passphrase"`
	Name     string `json:"name"`
}

func GetSSHDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %w", err)
	}
	return filepath.Join(home, ".ssh"), nil
}

func ListKeys() ([]SSHKey, error) {
	sshDir, err := GetSSHDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(sshDir)
	if err != nil {
		return nil, err
	}

	var privateNames = make(map[string]bool)
	var pubEntries []os.DirEntry

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".pub") {
			pubEntries = append(pubEntries, entry)
		} else if !strings.HasSuffix(name, ".pub") && !strings.HasSuffix(name, ".old") &&
			!strings.HasSuffix(name, ".bak") && name != "config" && name != "known_hosts" &&
			name != "authorized_keys" && name != "config.d" {
			privateNames[name] = true
		}
	}

	var keys []SSHKey

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		path := filepath.Join(sshDir, name)

		if strings.HasSuffix(name, ".pub") {
			baseName := strings.TrimSuffix(name, ".pub")
			if privateNames[baseName] {
				continue
			}

			pubContent, err := os.ReadFile(path)
			if err != nil {
				continue
			}

			keyType := detectKeyType(string(pubContent))
			fingerprint := ""
			if pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pubContent); err == nil {
				fingerprint = ssh.FingerprintSHA256(pubKey)
			}

			keys = append(keys, SSHKey{
				Name:          name,
				Path:          path,
				Type:          keyType,
				IsPrivate:     false,
				PubKeyContent: strings.TrimSpace(string(pubContent)),
				Fingerprint:   fingerprint,
			})
		} else if !strings.HasSuffix(name, ".pub") && !strings.HasSuffix(name, ".old") &&
			!strings.HasSuffix(name, ".bak") && name != "config" && name != "known_hosts" &&
			name != "authorized_keys" && name != "config.d" {
			pubPath := path + ".pub"
			keyType := "unknown"
			pubContent := ""
			fingerprint := ""

			if data, err := os.ReadFile(pubPath); err == nil {
				pubContent = strings.TrimSpace(string(data))
				keyType = detectKeyType(pubContent)
				if pubKey, _, _, _, err := ssh.ParseAuthorizedKey(data); err == nil {
					fingerprint = ssh.FingerprintSHA256(pubKey)
				}
			} else {
				keyType = detectPrivateKeyType(path)
			}

			keys = append(keys, SSHKey{
				Name:          name,
				Path:          path,
				Type:          keyType,
				IsPrivate:     true,
				PubKeyContent: pubContent,
				Fingerprint:   fingerprint,
			})
		}
	}

	return keys, nil
}

func detectKeyType(pubContent string) string {
	if strings.Contains(pubContent, "ssh-rsa") {
		return "RSA"
	} else if strings.Contains(pubContent, "ssh-ed25519") {
		return "ED25519"
	} else if strings.Contains(pubContent, "ecdsa-sha2-nistp") {
		return "ECDSA"
	} else if strings.Contains(pubContent, "ssh-dss") {
		return "DSA"
	}
	return "unknown"
}

func detectPrivateKeyType(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "unknown"
	}

	content := string(data)
	if strings.Contains(content, "RSA PRIVATE KEY") || strings.Contains(content, "OPENSSH PRIVATE KEY") {
		block, _ := pem.Decode(data)
		if block != nil {
			if strings.Contains(block.Type, "RSA") {
				return "RSA"
			}
			if privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
				switch privKey.(type) {
				case *ecdsa.PrivateKey:
					return "ECDSA"
				}
			}
		}
		return "RSA"
	} else if strings.Contains(content, "EC PRIVATE KEY") {
		return "ECDSA"
	} else if strings.Contains(content, "ED25519") {
		return "ED25519"
	}
	return "unknown"
}

func GenerateKey(req KeyGenRequest) (string, error) {
	sshDir, err := GetSSHDir()
	if err != nil {
		return "", err
	}

	if req.Name == "" {
		switch req.Type {
		case "rsa":
			req.Name = "id_rsa"
		case "ed25519":
			req.Name = "id_ed25519"
		case "ecdsa":
			req.Name = "id_ecdsa"
		default:
			req.Name = "id_ed25519"
		}
	}

	privatePath := filepath.Join(sshDir, req.Name)
	publicPath := privatePath + ".pub"

	if _, err := os.Stat(privatePath); err == nil {
		return "", fmt.Errorf("key file already exists: %s", privatePath)
	}

	var privateKey interface{}
	var pubKey ssh.PublicKey

	switch req.Type {
	case "rsa":
		bits := req.Bits
		if bits == 0 {
			bits = 4096
		}
		key, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			return "", fmt.Errorf("failed to generate RSA key: %w", err)
		}
		privateKey = key
		pubKey, err = ssh.NewPublicKey(&key.PublicKey)
		if err != nil {
			return "", fmt.Errorf("failed to create public key: %w", err)
		}

	case "ed25519":
		pub, priv, err := generateEd25519Key()
		if err != nil {
			return "", fmt.Errorf("failed to generate ED25519 key: %w", err)
		}
		privateKey = priv
		pubKey, err = ssh.NewPublicKey(ed25519.PublicKey(pub))
		if err != nil {
			return "", fmt.Errorf("failed to create public key: %w", err)
		}

	case "ecdsa":
		curve := elliptic.P256()
		if req.Bits == 384 {
			curve = elliptic.P384()
		} else if req.Bits == 521 {
			curve = elliptic.P521()
		}
		key, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return "", fmt.Errorf("failed to generate ECDSA key: %w", err)
		}
		privateKey = key
		pubKey, err = ssh.NewPublicKey(&key.PublicKey)
		if err != nil {
			return "", fmt.Errorf("failed to create public key: %w", err)
		}

	default:
		return "", fmt.Errorf("unsupported key type: %s", req.Type)
	}

	privBytes, err := marshalPrivateKey(privateKey, req.Passphrase)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private key: %w", err)
	}

	if err := os.WriteFile(privatePath, privBytes, 0600); err != nil {
		return "", fmt.Errorf("failed to write private key: %w", err)
	}

	pubComment := req.Email
	if pubComment == "" {
		pubComment = "generated-by-ssh-manager"
	}
	pubLine := string(ssh.MarshalAuthorizedKey(pubKey)) + pubComment + "\n"
	if err := os.WriteFile(publicPath, []byte(pubLine), 0644); err != nil {
		return "", fmt.Errorf("failed to write public key: %w", err)
	}

	return privatePath, nil
}

func generateEd25519Key() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return pub, priv, nil
}

func marshalPrivateKey(key interface{}, passphrase string) ([]byte, error) {
	if passphrase != "" {
		return marshalEncryptedPrivateKey(key, passphrase)
	}

	switch k := key.(type) {
	case *rsa.PrivateKey:
		privBytes := x509.MarshalPKCS1PrivateKey(k)
		block := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privBytes,
		}
		return pem.EncodeToMemory(block), nil
	case *ecdsa.PrivateKey:
		privBytes, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, err
		}
		block := &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: privBytes,
		}
		return pem.EncodeToMemory(block), nil
	case ed25519.PrivateKey:
		marshaled, err := ssh.MarshalPrivateKey(k, "")
		if err != nil {
			return nil, err
		}
		return pem.EncodeToMemory(marshaled), nil
	default:
		marshaled, err := ssh.MarshalPrivateKey(key, "")
		if err != nil {
			return nil, err
		}
		return pem.EncodeToMemory(marshaled), nil
	}
}

func marshalEncryptedPrivateKey(key interface{}, passphrase string) ([]byte, error) {
	marshaled, err := ssh.MarshalPrivateKey(key, passphrase)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(marshaled), nil
}

func DeleteKey(name string) error {
	sshDir, err := GetSSHDir()
	if err != nil {
		return err
	}

	privatePath := filepath.Join(sshDir, name)
	publicPath := privatePath + ".pub"

	if err := os.Remove(privatePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete private key: %w", err)
	}

	os.Remove(publicPath)
	return nil
}

func GetPubKeyContent(name string) (string, error) {
	sshDir, err := GetSSHDir()
	if err != nil {
		return "", err
	}

	pubPath := filepath.Join(sshDir, name+".pub")
	if _, err := os.Stat(pubPath); os.IsNotExist(err) {
		return "", fmt.Errorf("public key not found: %s", pubPath)
	}

	content, err := os.ReadFile(pubPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

func GetPrivKeyContent(name string) (string, error) {
	sshDir, err := GetSSHDir()
	if err != nil {
		return "", err
	}

	privPath := filepath.Join(sshDir, name)
	if _, err := os.Stat(privPath); os.IsNotExist(err) {
		return "", fmt.Errorf("private key not found: %s", privPath)
	}

	content, err := os.ReadFile(privPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

func CopyPubKeyToClipboard(name string) (string, error) {
	return GetPubKeyContent(name)
}
