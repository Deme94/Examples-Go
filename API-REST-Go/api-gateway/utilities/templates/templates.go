package templates

const CONFIRM_EMAIL string = `<body>
<img id="logo" src="https://revistabyte.es/wp-content/uploads/2022/07/que-es-un-desarrollador-de-go-y-como-convertirse-en-uno.jpg"/>
<div id="main">
 <h2>Confirm Your Email Address</h2>
 <p>Tap the button below to confirm your email address. 
 If you didn't create an account with Paste, you can safely delete this email.</p>
 <button href="FRONT_DOMAINCONFIRM_EMAIL_ROUTECONFIRM_EMAIL_TOKEN">Confirm Email</button>
 <p>If that doesn't work, click on the link below:</p>
 <a href="FRONT_DOMAINCONFIRM_EMAIL_ROUTECONFIRM_EMAIL_TOKEN">Click here</a>
 <div id="cheers">
   <p>Cheers,</p>
   <p>COMPANY_OWNER</p>
 </div>
</div>
<div id="ending">
 <p>You received this email because we received a request for Signup 
   for your account. If you didn't request Signup you can safely delete 
   this email.</p>

 <p>To stop receiving these emails, you can <a href="FRONT_DOMAIN">unsubscribe</a> at any time.</p>

 <p>COMPANY_NAME</p>
</div>
</body>

<style>
body {
 background-color: gray;
 display: flex;
 justify-content: center;
 padding: 30px;
 flex-direction: column;
 place-items: center;
 gap: 30px;
 background-color: #e9ecef;
}
#logo {
width: 100px;
}
#main {
 padding: 30px;
 text-align: center;
 width: 50%;
 background-color: white;
 border-top: 5px solid #55551c;
}
#cheers {
 padding-top: 30px;
 text-align: left;
 line-height: 10%;
}
#ending {
 width: 50%;
 font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif;
 font-size: 14px;
 line-height: 20px;
 color: #666;
 text-align: center;
}
button {
 display: inline-block;
 padding: 16px 36px;
 font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif;
 font-size: 16px;
 color: #ffffff;
 text-decoration: none;
 border-radius: 6px;
 background-color: rgb(26, 130, 226);
 border: 0;
}
</style>`
