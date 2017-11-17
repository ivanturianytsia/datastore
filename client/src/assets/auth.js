class Auth {
  constructor (ctx) {
    this.ctx = ctx
  }
  post (url, email, password) {
    if (email === '') {
      throw new Error('No email provided')
    }
    if (password === '') {
      throw new Error('No password provided')
    }
    const that = this
    return this.ctx.$http.post(url, {
      email,
      password
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
  Login (email, password) {
    return this.post('/auth/login', email, password)
  }
  Register (email, password) {
    return this.post('/auth/register', email, password)
  }
  GetUser () {
    if (this.token === '') {
      throw new Error('No token provided')
    }
    return this.ctx.$http.get('/auth/user', {
      headers: {
        'Authentication': 'Bearer ' + this.token
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
