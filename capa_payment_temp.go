package helpers

const CapaPaymentTemplate = `<div style="font-family:Helvetica Light,Helvetica,Arial,sans-serif;margin:0;padding:0; width:100%" bgcolor="#eeeeee">  

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
<td align="center" style="color:#000000;font-family:Arial,Helvetica,sans-serif;font-size:18px;font-weight:bold;padding:0px;padding-top: 40px;" >
<img src="{{AUTOMATED_EMAIL_LOGO}}" alt=" logo" width="250" height="45" />
</td>
</tr>

<tr>
<td align="center" style="color:#000000;font-family:Arial,Helvetica,sans-serif;font-size:25px;font-weight:bold;padding:0px;padding-top: 40px;" >
Payment Received
</td>
</tr>
<tr>
<td align="left" style="color:#000000;font-family:Arial,Helvetica,sans-serif;font-size:15px;font-weight:normal;line-height:22px;padding:30px 5% 0px" >
<p style="line-height:1.3em">
<h3 style="margin-bottom:.6em">Dear {{fullname}}</h3>
Your payments of <b>${{amount}}</b> for extra <b>{{capacity}}</b> agent capacity was received on <b>{{date_received}}.</b>
</p>


  
<h4 style="color:#41474e;line-height:23px;font-size:19px;margin-bottom:8px;font-weight:normal">Payment details</h4>
<center> 
<table width="100%" style="color:#5e6670;font-family:Helvetica,Arial,sans-serif;font-size:15px;line-height:1.25em;border-collapse:collapse;border-top:1px solid #e4e6e9">
<tbody>
<tr>
<td width="33%" style="border-bottom:1px solid #e4e6e9;text-align:left;vertical-align:top;width:33%;color:black;padding:10px 5px;background:#eef6e5;font-weight:bold">
<span style="color:black">Amount Received</span></td>
<th align="left" style="border-bottom:1px solid #e4e6e9;vertical-align:top;color:black;padding:10px 5px;background:#eef6e5;font-weight:bold;text-align:right">
${{amount}} <br/>
</th>
</tr>  

<tr>
<td style="border-bottom:1px solid #e4e6e9;color:#7a828c;padding:10px 3px;text-align:left;vertical-align:top;width:33%">Number Of Capacities:</td>
<td style="border-bottom:1px solid #e4e6e9;color:#7a828c;padding:10px 3px;text-align:left;vertical-align:top;width:33%">
<b>{{capacity}}</b>
</td>
</tr>

<tr>
<td style="border-bottom:1px solid #e4e6e9;color:#7a828c;padding:10px 3px;text-align:left;vertical-align:top;width:33%">Refrence ID:</td>
<td style="border-bottom:1px solid #e4e6e9;color:#7a828c;padding:10px 3px;text-align:left;vertical-align:top;width:33%">
<a style="color: #0098E1;">{{refrence_id}}</a>
</td>
</tr>

<tr>
<td style="border-bottom:1px solid #e4e6e9;color:#7a828c;padding:10px 3px;text-align:left;vertical-align:top;width:33%">Transaction ID:</td>
<td style="border-bottom:1px solid #e4e6e9;color:#7a828c;padding:10px 3px;text-align:left;vertical-align:top;width:33%">
<a style="color: #0098E1;">{{transaction_id}}</a>
</td>
</tr>
  
</tbody>
</table> 
</center>
 
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

<tbody>
<tr>
<td bgcolor="#2ccae7" height="46" align="center" style="border-radius:2px;">
<a href="{{crm_url}}/login.php" style="color:#ffffff;display:inline-block;font-family:\'Helvetica Neue\',arial;font-size:17px;font-weight:bold;line-height:46px;min-width:280px;max-width:280px;text-align:center;text-decoration:none">
Login To {{app_name}}
</a>
</td>
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
<td style="color:#535353;font-size:10px;line-height:16px;padding-bottom:20px;padding-left: 15px;padding-right: 15px;" align="center">
<span style="font-size:12px"><span style="font-family:arial,helvetica neue,helvetica,sans-serif">
<a href="{{base_url}}/privacy.php">Privacy Policy</a> 
| <a href="{{base_url}}/terms.php">Terms of Use</a>&nbsp;<br>
This email was sent to you by <a href="{{base_url}}">{{app_name}}</a>. <br/>You are receiving this email because you signed up for <a href="{{base_url}}">{{app_name}}</a>.&nbsp;<br>
<br>
<span>{{app_name}}</span> | {{comp_address}}</span></span></td>
</tr>
</tbody>
</table>
<br>
<br>
&nbsp;
</td>
</tr>
</tbody></table>

</td>
</tr>
</tbody>
</table>
</div>`
