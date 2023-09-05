const mongoose = require('mongoose');
require('dotenv').config();
mongoose.set('strictQuery', true);

const uri = process.env.URI

const connectDB = () =>{
    return mongoose.connect(uri,{ useNewUrlParser: true, useUnifiedTopology: true },()=>{
        console.log("Connection successfull")
    });
}

module.exports = connectDB