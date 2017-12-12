const domain = ''

class Auth {
  constructor (ctx) {
    this.ctx = ctx
  }
  post (url, credentials) {
    return this.ctx.$http.post(url, credentials, {
      headers: {
        'Content-Type': 'application/json'
      }
    }).then(response => {
      let obj = JSON.parse(response.body)
      if (obj.token) {
        this.token = obj.token
      }
      return obj
    })
  }
  Login (email, password) {
    if (email === '') {
      throw new Error('No email provided')
    }
    if (password === '') {
      throw new Error('No password provided')
    }
    return this.post(domain + '/auth/login', {
      email,
      password
    })
  }
  Register (email, password, phonenumber) {
    if (email === '') {
      throw new Error('No email provided')
    }
    if (password === '') {
      throw new Error('No password provided')
    }
    if (phonenumber === '') {
      throw new Error('No phonenumber provided')
    }
    return this.post(domain + '/auth/register', {
      email,
      password,
      phonenumber
    })
  }
  PostEmailCode (email, code) {
    if (email === '') {
      throw new Error('No email provided')
    }
    if (code === '') {
      throw new Error('No code provided')
    }
    const that = this
    return this.ctx.$http.post('/auth/code', {
      email,
      code
    }, {
      headers: {
        'Content-Type': 'application/json'
      }
    }).then(response => {
      let obj = JSON.parse(response.body)
      if (obj.token) {
        that.token = obj.token
      }
      return obj
    })
  }
  GetUser () {
    if (this.token === '') {
      throw new Error('No token provided')
    }
    return this.ctx.$http.get(domain + '/auth/user', {
      headers: {
        'Authorization': 'Bearer ' + this.token
      }
    })
    .then(response => {
      return JSON.parse(response.body)
    })
  }
  PutUser (updates) {
    if (this.token === '') {
      throw new Error('No token provided')
    }
    return this.ctx.$http.put(domain + '/user', updates, {
      headers: {
        'Authorization': `Bearer ${this.token}`,
        'Content-Type': 'application/json'
      }
    })
    .then(response => {
      return JSON.parse(response.body)
    })
  }
  GetUsersByEmail (email) {
    if (this.token === '') {
      throw new Error('No token provided')
    }
    return this.ctx.$http.get(`${domain}/users?email=${email}`, {
      headers: {
        'Authorization': 'Bearer ' + this.token
      }
    })
    .then(response => {
      return JSON.parse(response.body)
    })
  }
  Logout () {
    this.token = ''
  }
  IsLogged () {
    return !!this.token
  }
  get token () {
    return localStorage.token
  }
  set token (val) {
    localStorage.token = val
  }
}

// Login
// Register
// Logout
// IsLogged
// User

export default Auth
