package helpers

const NewDealNotificationTemplate = `Hi {{agent_name}}, <br>
You've received new {{type_out}} on {{APP_NAME}}. You now have <b>{{total_of}}</b> you haven't checked yet. 
<br />
<br />
<br />
<center>
<center>
<table width="800" style="color:#5e6670;font-family:Helvetica,Arial,sans-serif;font-size:15px;line-height:1.25em;margin:0 auto">
<tbody><tr>
<td align="center">
<a href="{{crm_url}}/login.php" style="text-decoration:none;color:#ffffff;background:#009EEA;border-radius:50px;display:inline-block;font-size:18px;
font-weight:bold;margin-bottom:5px;padding:12px 80px;white-space:nowrap" target="_blank">
Login to your dashboard <img src="{{crm_url}}/assets/img/arr.png" alt="" width="11" height="17" style="border:0;outline:none;
text-decoration:none;margin-left:3px;vertical-align:-2px">
</a>
</td>
</tr>
</tbody></table>
</center> 
</center>`
