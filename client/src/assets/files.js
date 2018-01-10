const domain = ''

class Files {
  constructor (ctx) {
    this.ctx = ctx
  }
  GetFiles () {
    if (this.token === '') {
      throw new Error('No token provided')
    }
    return this.ctx.$http.get(domain + '/files', {
      headers: {
        'Authorization': 'Bearer ' + this.token
      }
    })
    .then(response => {
      let obj = {}
      try {
        obj = JSON.parse(response.body)
      } catch (err) {
        obj = response.body
      }
      return obj
    })
  }
  DeleteFile (filename) {
    if (this.token === '') {
      throw new Error('No token provided')
    }
    return this.ctx.$http.delete(`${domain}/files/${filename}`, {
      headers: {
        'Authorization': 'Bearer ' + this.token
      }
    })
  }
  UpdateFile (fileid, updates) {
    if (this.token === '') {
      throw new Error('No token provided')
    }
    return this.ctx.$http.put(`${domain}/file/${fileid}`, updates, {
      headers: {
        'Authorization': 'Bearer ' + this.token,
        'Content-Type': 'application/json'
      }
    })
  }
  UploadFiles (files) {
    if (this.token === '') {
      throw new Error('No token provided')
    }

    var data = new FormData()
    data.append('file', files[0])

    return this.ctx.$http.post(`${domain}/upload`, data, {
      headers: {
        'Authorization': 'Bearer ' + this.token,
        'Content-Type': 'multipart/form-data'
      }
    })
    .then(response => {
      let obj = {}
      try {
        obj = JSON.parse(response.body)
      } catch (err) {
        obj = response.body
      }
      return obj
    })
  }
  get token () {
    return localStorage.token
  }
  set token (val) {
    localStorage.token = val
  }
}

export default Files
