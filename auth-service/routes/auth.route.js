const express = require('express')
const router = express.Router()
const { logOut, signin_post, signup_post} = require('../controllers/auth.controllers')
const verifyToken = require('../middlewares/authJWT')

router.route('/api/auth/signup').post(signup_post)
router.route('/api/auth/signin').post(signin_post)
router.route('/api/auth/logout').get(verifyToken,logOut)

module.exports = router
