const express = require('express')
const db = require('./db/conn')
const bodyParser = require('body-parser')
const app = express()
const auth = require("./routes/auth.route")
const cookieParser = require('cookie-parser')
const port = process.env.PORT

app.use(cookieParser())
app.use(express.json());
app.use(express.urlencoded({extended:false}));
app.use(bodyParser.json());

app.use('/',auth)

const start = async () => {
    try {
        await db();
        app.listen(port,()=>{
            console.log("Server Listening on port",port); 
         })
    } catch (error){
        console.log(error);
    }
}

start();
