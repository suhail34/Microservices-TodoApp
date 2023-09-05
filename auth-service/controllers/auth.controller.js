const jwt = require('jsonwebtoken');
const User = require('../models/user.models');
const cookie = require('cookie-parser');

const handleErrors = (err) => {
  let errors = {fullName:'',email:'',role:'',password:''};

  if(err.code===11000){
    errors.email="Email already exists";
    return errors;
  }

  if(err.message === 'Invalid Email'){
    errors.email = "Email is Incorrect";
    return errors;
  }

  if(err.message === 'Invalid Password'){
    errors.password = "Incorrect Password";
    return errors;
  }

  if(err.message === 'Not a valid Token'){
    return {err:err.message}
  }

  if(err.message.includes('User validation failed')) {
    Object.values(err.errors).forEach(({properties}) => {
      errors[properties.path] = properties.message;
    });
  }
  return errors;
}

const maxAge = 1*24*60*60;
const createToken = (id) => {
  return jwt.sign({id},process.env.API_SECRET,{expiresIn:maxAge});
}

/**
 * Function for User Registration
 * @param {*} req 
 * @param {*} res 
 */
const signup_post = async (req, res) => {
  const {fullName,email,role,password} = req.body
  try {
    const user = await User.create({fullName,email,role,password})
    console.log(user)
    const token = createToken(user._id)
    res.cookie('jwt',token,{expiresIn:1000*maxAge,httpOnly:true})
    res.status(201).send({user:user._id})
  }catch (err) {
    const errors = handleErrors(err)
    res.status(400).json({errors})
  }
}

/**
 * Function for user login
 * @param {*} req 
 * @param {*} res 
 */
const signin_post = async (req, res) => {
  const {email, password} = req.body;
  console.log(email,password)
   try {
    const user = await User.login(email,password)
    const token = createToken(user._id)
    res.cookie('jwt',token,{expiresIn:1000*maxAge,httpOnly:true})
    res.status(200).json({user:user._id})
   } catch (err){
    const errors = handleErrors(err)
    res.status(400).json({error:errors})
   }
}

// const signin = (req, res) => {
//   User.findOne({
//     email: req.body.email
//   }).exec((err, user) => {
//     if (err) {
//       res.status(500).send({
//         message: err
//       })
//       return;
//     }
//     if (!user) {
//       return res.status(404).send({
//         message: "User Not Found. Please Register"
//       })
//     }
//     const passwordValidity = bcrypt.compareSync(
//       req.body.password,
//       user.password
//     )

//     if (!passwordValidity) {
//       return res.status(401).send({
//         accessToken: null,
//         message: "Not a valid Password"
//       })
//     }

//     const token = jwt.sign({
//       id: user.id
//     }, process.env.API_SECRET, {
//       expiresIn: 86400
//     })


//     res.status(200)
//       .send({
//         user: {
//           id: user._id,
//           email: user.email,
//           fullName: user.fullName,
//         },
//         message: "Login successfull",
//         accessToken: token,
//       });
//   })
// }

/**
 * Function for User logout
 * @param {*} req 
 * @param {*} res 
 */
const logOut = (req, res) => {
  res.cookie('jwt','',{maxAge:1});
  res.redirect('/api/auth/signin');
}

module.exports = { signup_post, logOut, signin_post }
