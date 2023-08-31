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
    * @param {string} pattern - string to search torrents for
    * @return {Promise<number>} id - id of the search
    */
   async initSearch(pattern) {
      const res = await this.#post("api/v2/search/start", { pattern: pattern, plugins: "enabled", category: "all" })
      return res.data.id
   }

   /**
    * Stops torrrent search
    * @param {number} id - search id
    */
   async stopSearch(id) {
      await this.#post("api/v2/search/stop", { id: id })
   }

   /**
    * Deletes torrrent search
    * @param {number} id - search id
    */
   async deleteSearch(id) {
      await this.#post("api/v2/search/delete", { id: id })
   }

   /**
    * Get torrrent search results
    * @param {number} id - search id
    * @param {number} limit - search results limit, 0 => no limit
    * @return {Promise<object>} searchResults - object containing search results, status and count of results
    */
   async getSearchResults(id, limit = 0) {
      const res = await this.#get("api/v2/search/results", { id: id, limit: limit })
      return res.data
   }

   /**
    * Add torrrent to download
    * @param {string} urls - URLs separated with newlines
    * @param {string} savepath - Download folder
    * @param {string} root_folder - Create the root folder. Possible values are true, false, unset (default)
    * @param {string} sequentialDownload - Enable sequential download. Possible values are true, false (default)
    * @param {string} firstLastPiecePrio - Prioritize download first last piece. Possible values are true, false (default)
    * @param {string} rename - Rename torrent
    * @return {Promise<boolean>} success - true/false
    */
   async addTorrentDownload(urls, savepath, root_folder = "true", sequentialDownload = "true", firstLastPiecePrio = "true", rename) {
      const res = await this.#post("api/v2/torrents/add", { urls: urls, savepath: savepath, root_folder: root_folder, sequentialDownload: sequentialDownload, firstLastPiecePrio: firstLastPiecePrio, rename: rename })
      if (res.data == "Ok.") return true
      return false
   }

   /**
    * Get torrrent Hash
    * @param {string} name - Name of the torrent to fetch hash of
    * @param {string} filter - Filter torrent list by state. Allowed state filters: all, downloading, seeding, completed, paused, active, inactive, resumed, stalled, stalled_uploading, stalled_downloading, errored
    * @return {Promise<string>} hash - Hash of the torrent
    */
   async getTorrentHash(name, filter = "downloading") {
      const res = await this.#get("api/v2/torrents/info", { filter: filter })
      const torrents = res.data
      return torrents.filter((torrent) => torrent.name === name)[0].hash
   }

   /**
    * Delete torrrent
    * @param {string} hashes - Hash of torrent to delete
    * @param {string} deleteFiles - If set to true, the downloaded data will also be deleted, otherwise has no effect.
    */
   async deleteTorrent(hashes, deleteFiles = "false") {
      await this.#post("api/v2/torrents/delete", { hashes: hashes, deleteFiles: deleteFiles })
   }
}
