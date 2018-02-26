import auth0 from 'auth0-js'

import {
  AUTH_CONFIG
} from './auth0-config'

var auth0Instance = new auth0.WebAuth({
  domain: AUTH_CONFIG.domain,
  clientID: AUTH_CONFIG.clientId,
  redirectUri: AUTH_CONFIG.callbackUrl,
  audience: AUTH_CONFIG.audience,
  responseType: 'token id_token',
  scope: 'openid profile email'
})

var LoginService = {
  methods: {
    auth0Login() {
      // show auth0 authorize login
      auth0Instance.authorize()
    },
    auth0Logout() {
      // clear local storage
      localStorage.removeItem('access_token')
      localStorage.removeItem('id_token')
      localStorage.removeItem('expires_at')
    },
    auth0Handler(callback) {
      // parse hash url
      auth0Instance.parseHash((err, authResult) => {
        if (err) {
          return console.log(err)
        }

        // get user information
        auth0Instance.client.userInfo(authResult.accessToken, function (err, user) {
          if (err) {
            return console.log(err)
          }

          // save user information
          localStorage.setItem('logged_user', JSON.stringify(user))

          // redirect to dashboard
          if (authResult && authResult.accessToken && authResult.idToken) {
            localStorage.setItem('access_token', authResult.accessToken)
            localStorage.setItem('id_token', authResult.idToken)
            localStorage.setItem('expires_at', JSON.stringify(
              authResult.expiresIn * 1000 + new Date().getTime()
            ))
            callback('dashboard')
          } else if (err) {
            callback('login')
            console.log(err)
          }
        })
      })
    },
    auth0CheckAuth() {
      // check if token is expired
      return new Date().getTime() < JSON.parse(localStorage.getItem('expires_at'))
    }
  }
};
export default LoginService;
