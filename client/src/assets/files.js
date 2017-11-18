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
  get token () {
    return localStorage.token
  }
  set token (val) {
    localStorage.token = val
  }
}

export default Files
