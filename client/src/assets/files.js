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
      return JSON.parse(response.body)
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
      return JSON.parse(response.body)
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
