const axios = require("axios")
require("dotenv").config()

module.exports.qBittorent = class {
   #cookie
   #username
   #password
   performReq
   constructor(username, password, host) {
      this.#username = username
      this.#password = password
      this.performReq = axios.create({ baseURL: host, headers: { Referer: host, Host: host, "Content-Type": "application/x-www-form-urlencoded" } })
   }

   async post(url, data) {
      const res = await this.performReq.post(url, data, { headers: { Cookie: this.#cookie } })
   }

   async get(url, params = {}) {
      const res = await this.performReq.get(url, { headers: { Cookie: this.#cookie }, params })
      console.log(res)
   }

   async getNewAuthSession() {
      const res = await this.performReq.post("api/v2/auth/login", { username: this.#username, password: this.#password })
      if (res.headers["set-cookie"]) {
         this.#cookie = res.headers["set-cookie"]
         return true
      }
      return false
   }
}
