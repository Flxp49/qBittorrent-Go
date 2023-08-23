const axios = require("axios")

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

   /** Fetches new auth cookie
    */
   async getNewAuthSession() {
      try {
         const res = await this.#post("api/v2/auth/login", { username: this.#username, password: this.#password })
         if (res.headers["set-cookie"]) {
            this.#cookie = res.headers["set-cookie"]
            return
         } else {
            throw new Error(`Login Failed. Username: ${this.#username} Message: ${res.data}`)
         }
      } catch (error) {
         throw new Error(`Login Failed. Username: ${this.#username} Error: ${error}`)
      }
   }

   async #post(url, data) {
      return await this.performReq.post(url, data, { headers: { Cookie: this.#cookie } })
   }

   async #get(url, params = {}) {
      return await this.performReq.get(url, { headers: { Cookie: this.#cookie }, params })
   }

   /**
    * Starts torrrent search
    * @param {string} searchpattern - string to search torrents for
    * @return {number} - id of the search
    */
   async initSearch(searchpattern) {
      const res = await this.#post("api/v2/search/start", { pattern: searchpattern, plugins: "enabled", category: "all" })
      if (res.data.id) {
         return res.data.id
      } else {
         return false
      }
   }

   /**
    * Stops torrrent search
    * @param {number} sid - search id
    */
   async stopSearch(sid) {
      await this.#post("api/v2/search/stop", { id: sid })
      return
   }

   /**
    * Deletes torrrent search
    * @param {number} sid - search id
    */
   async stopSearch(sid) {
      await this.#post("api/v2/search/delete", { id: sid })
      return
   }

   /**
    * Fetch torrrent search results
    * @param {number} sid - search id
    * @param {number} slimit - search results limit, 0 => no limit
    * @return {object} - object containing search results, status and count of results
    */
   async fetchSearchResults(sid, slimit = 0) {
      const res = await this.#get("api/v2/search/results", { id: sid, limit: slimit })
      return res
   }

   /**
    * Add torrrent to download
    * @param {string} url - URLs separated with newlines
    * @param {string} path - Download folder
    * @param {string} rootFolder - Create the root folder. Possible values are true, false, unset (default)
    * @param {string} seqDownload - Enable sequential download. Possible values are true, false (default)
    * @param {string} firstLastPriority - Prioritize download first last piece. Possible values are true, false (default)
    */
   async addTorrentDownload(url, path, rootFolder = "true", seqDownload = "true", firstLastPriority = "true") {
      const res = await this.#post("api/v2/torrents/add", { urls: url, savepath: path, root_folder: rootFolder, sequentialDownload: seqDownload, firstLastPiecePrio: firstLastPriority })
      return
   }
}
