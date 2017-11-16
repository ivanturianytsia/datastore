class AuthApi {
  constructor () {

  }
  send (url, email, password) {
    if (email === "") {
      throw new Error('No email provided')
    }
    if (password === "") {
      throw new Error('No password provided')
    }
    return $.post({
      url: url,
      contentType: "application/json",
      data: JSON.stringify({
        email,
        password
      })
    })
  }
  Login (email, password) {
    return this.send('/auth/login', email, password)
  }
  Register (email, password) {
    return this.send('/auth/register', email, password)
  }
  GetUser () {

  }
}


$(() => {
  const auth = new AuthApi()

  $('#register-form').submit(event => {
    event.preventDefault()
    auth.Register($('#email').val(), $('#password').val())
      .then(response => {
        console.log(response)
      })
      .catch(err => {
        console.log(err)
      })
  })

  $('#login-form').submit(event => {
    event.preventDefault()

    auth.Login($('#email').val(), $('#password').val())
      .then(response => {
        console.log(response)
      })
      .catch(err => {
        console.log(err)
      })
  })
})
