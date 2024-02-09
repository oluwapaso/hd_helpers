package helpers

const VerifyEmailTemplate = `<div style="font-family:Helvetica Light,Helvetica,Arial,sans-serif;margin:0;padding:0; width:100%" bgcolor="#eeeeee">  
          
<table border="0" cellpadding="0" cellspacing="0" width="100%" style="border-collapse:collapse">
<tbody><tr>
<td bgcolor="#eeeeee" align="center" style="padding:25px" >

<table bgcolor="#ffffff" border="0" cellpadding="0" cellspacing="0" width="100%" style="border-collapse:collapse;max-width:600px" >
<tbody><tr>
<td>

<table width="100%" border="0" cellspacing="0" cellpadding="0" style="border-collapse:collapse">
<tbody>
<tr>
<td>

<table width="100%" border="0" cellspacing="0" cellpadding="0" style="border-collapse:collapse">
<tbody>

<tr>
<td align="center" style="color:#000000;font-family:Arial,Helvetica,sans-serif;font-size:25px;font-weight:bold;padding:0px;padding-top: 40px;" >
Welcome to {{app_name}}
</td>
</tr>
<tr>
<td align="center" style="color:#000000;font-family:Arial,Helvetica,sans-serif;font-size:15px;font-weight:normal;line-height:22px;padding:30px 5% 0px" >
We\'re so excited to show you over the next two weeks (and hopefully beyond) how you can improve your auto transport business. 
And we are confident we can accomplish a lot together. 
<br /><br />
Verify your email address to complete registration
</td>
</tr>
        
<tr>
<td width="100%" align="center" valign="top" bgcolor="#ffffff" height="20"></td>
</tr>
 

<tr>
<td width="100%" align="center" valign="top" bgcolor="#ffffff" height="20"></td>
</tr>



<tr>
<td width="100%" align="center" valign="top" bgcolor="#ffffff" height="1" style="padding:0px 30px">
<table cellpadding="0" cellspacing="0" width="30%" style="border-collapse:collapse">
<tbody><tr>
<td style="border-top-color:#eeeeee;border-top-style:solid;border-top-width:1px;padding:0px 30px"></td>
<td>
</td>
</tr>
</tbody></table>
</td>
</tr>

 
</tbody></table>
</td>
</tr>
 


<tr>
<td>

<table width="100%" border="0" cellspacing="0" cellpadding="0" style="border-collapse:collapse">

<tbody>  
<tr>
<td width="100%" align="center" valign="top" bgcolor="#ffffff" height="20"></td>
</tr>



<tr>
<td width="100%" align="center" valign="top" bgcolor="#ffffff" height="1" style="padding:0px 30px">
<table cellpadding="0" cellspacing="0" width="300" height="46" style="border-collapse:collapse">

<tbody><tr>
<td bgcolor="#2ccae7" height="46" align="center" style="border-radius:2px;">
<a href="{{base_url}}/select-account-type.php?email={{email}}&token={{token}}" 
style="color:#ffffff;display:inline-block;font-family:\'Helvetica Neue\',arial;font-size:17px;font-weight:bold;line-height:46px;min-width:280px;max-width:280px;text-align:center;text-decoration:none">
Verify &amp; Complete Registration</a></td>
</tr>
</tbody>
</table>
</td>
</tr>
 

<tr>
<td width="100%" align="center" valign="top" bgcolor="#ffffff" height="20"></td>
</tr>


<tr>
<td width="100%" align="center" valign="top" bgcolor="#ffffff" height="20"></td>
</tr>
</tbody></table>
</td>
</tr>
 

<tr>
<td align="center" valign="top" style="font-size:0" >

</td>
</tr>
</tbody>
</table>
</td>
</tr>
</tbody>
</table>
</td>
</tr> 

<tr>
<td bgcolor="#eeeeee" align="center" style="padding:20px 0px">


<table width="100%" border="0" cellspacing="0" cellpadding="0" align="center" style="border-collapse:collapse;max-width:600px" >
<tbody><tr>
<td align="center" style="color:#818181;font-family:\'Helvetica Light\',\'Helvetica\',Arial,sans-serif;font-size:12px;line-height:1.5;padding-top:5px">

<table style="border-collapse:collapse;text-align:center;width:100%;">
<tbody>
<tr>
<td style="color:#535353;font-size:10px;line-height:16px;padding-bottom:20px;padding-left: 15px;padding-right: 15px;" align="center"><span style="font-size:12px"><span style="font-family:arial,helvetica neue,helvetica,sans-serif">
<a href="{{base_url}}/privacy.php">Privacy Policy</a> | <a href="{{base_url}}/terms.php">Terms of Use</a>&nbsp;<br>
<span{{app_name}}</span> {{comp_address}}</span></span></td>
</tr>
</tbody>
</table>
<br>
<br>
&nbsp;
</td>
</tr>
</tbody>
</table>

</td>
</tr>
</tbody>
</table>
</div>
</div>`
