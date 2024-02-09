package helpers

const TicketResponseTemplate = `<center>
    <table style="color: #5e6670; font-family: Helvetica,Arial,sans-serif; font-size: 15px; line-height: 1.25em; background-color: white; padding: 20px 0; text-align: left; padding-bottom: 20px;" width="600" align="center">
    <tbody>
    <tr>
    <td align="center"><img class="CToWUd" style="border: 0; outline: none; width: auto; max-width: 300px; height: auto; text-decoration: none;" src="" alt="" /></td>
    </tr>
    </tbody>
    </table>
    <table style="color: #5e6670; font-family: Helvetica,Arial,sans-serif; font-size: 15px; line-height: 1.25em; background-color: white; padding: 0; text-align: left; padding-bottom: 20px;" width="600" align="center">
    <tbody>
    <tr>
    <td><hr style="background: #e4e6e9; color: #e4e6e9; font-size: 1px; height: 1px; border: 0;" /></td>
    </tr>
    <tr>
    <td>
    <h3 style="width: 100%; text-align: left; font-weight: bold; padding-top:30px;">Hi {{firstname}}</h3>
    <p>Thank you for contacting <b>{{app_name}}</b></p>
    
    <p>This is just a quick note to let you know we've received your message, and we will respond as soon as we can.</p>
    <p>Please use <b>{{ticket_number}}</b> to track your issue</p>
    
    <center>
    <table style="color:#5e6670;font-family:Helvetica,Arial,sans-serif;font-size:15px;line-height:1.25em;margin:0 auto">
    <tbody>
    <tr>
    <td align="center">
    <a href="{{crm_url}}/#/ticket-details/{{ticket_id}}" style="text-decoration:none;color:#ffffff;background:#009EEA;border-radius:50px;display:inline-block;font-size:18px;font-weight:bold;margin-bottom:5px;padding:12px 80px;white-space:nowrap" target="_blank">
    Track Issue Details <img src="https://hauingdesk.com/crm/assets/img/arr.png" alt="" width="11" height="17" style="border:0;outline:none;text-decoration:none;margin-left:3px;vertical-align:-2px" class="CToWUd">
    </a>
    </td>
    </tr>
    </tbody></table>
    </center> 
    
    <center>
    <p style="line-height: 1.3em; color: black;"><strong>Questions or concerns? Please contact us immediately.</strong></p>
    </center>&nbsp;
    
    </tbody>
    </table>
    </center>`
