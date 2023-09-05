const mongoose = require('mongoose');
const bcrypt = require('bcrypt');

const userSchema = new mongoose.Schema({
    fullName: {
        type: String,
        required: [true, "Full name not provided"],
        minLength:[6,"Enter Your Full Name"]
      },
      email: {
        type: String,
        unique: true,
        lowercase: true,
        trim: true,
        required: [true, "Email not provided"],
        validate: {
          validator: function (v) {
            return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v);
          },
          message: '{VALUE} is not a valid email!'
        }
    
      },
      role: {
        type: String,
        enum: ["normal", "admin"],
        required: [true, "Please specify user role"]
      },
      password: {
        type: String,
        required: [true,"Password Not Provided"],
        minLength: [3,"Minimum Password Length is 8"]
      },
      created: {
        type: Date,
        default: Date.now
      }
})

userSchema.pre('save',async function (next) {
  const salt = await bcrypt.genSalt();
  this.password = await bcrypt.hash(this.password,salt);
  // console.log("pre__",this)
  next()
})

userSchema.statics.login = async function(email,password) {
  const user = await this.findOne({email});
  if(user){
    const auth = await bcrypt.compare(password,user.password);
    if(auth){
      return user;
    }
    throw Error('Invalid Password');
  }
  throw Error('Invalid Email');
}

/**
 * User Registration Model
 * @param {String} fullName
 * @param {String} email
 * @param {String} role
 * @param {String} password
 */
module.exports = mongoose.model('User',userSchema)
