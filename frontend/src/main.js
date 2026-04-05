import { createApp } from 'vue'
import App from './App.vue'

function waitForWails() {
  return new Promise((resolve) => {
    if (window.go) {
      resolve()
      return
    }
    const check = setInterval(() => {
      if (window.go) {
        clearInterval(check)
        resolve()
      }
    }, 50)
  })
}

waitForWails().then(() => {
  createApp(App).mount('#app')
})
