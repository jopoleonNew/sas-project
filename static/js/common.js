var currentUrl = window.location.href
$(document).ready(function () {

		$('#formcontainer').keypress(function(e) {
			if (e.keyCode == 13) {
				$('#loginbutton').click();
			}
				
		});
	

	$("#gotosignup").click(function (e) {
		e.preventDefault();
		window.location = currentUrl + "signup";
	});
	$("#gotologin").click(function (e) {
		e.preventDefault();
		//window.location = currentUrl.replace(" ", "login");
		window.location = currentUrl + "login";
	});
$("#formcontainer").change(function () {
	$('button[type=submit]').prop('disabled',false);
});
	console.log(currentUrl)
	$("#loginbutton").click(function () {
		console.log($("#userName").val())
		console.log($("#password").val())


		window.userName = $("#userName").val();
		window.password = $("#password").val();

		$('button[type=submit]').prop('disabled',true);
		$.ajax({
			data: {
				"username": window.userName,
				"password": window.password,

			},
			//dataType: "json",
			type: "POST",


			url: currentUrl + "loginsubmit",


			success: function (data) {
				//location.reload();
				$('#loginresult').empty()
				if (data.indexOf("Success") >= 0) {
					console.log("succsessful login")
					$('#loginresult').append("<p style='color: #5CB85C; height: 30px;font-size: 18px'>"+data+"</p>");
					
					setTimeout(function () {
						location.reload();
					}, 1000);
				} else {
					$('#loginresult').append(
					"<p style='color: #DE5246; height: 30px;font-size: 18px'>"+data+"</p>"
				);
				}
				
					//console.log("Data sent: ", data)

				console.log("Data resieved: ", data)

			},
			error: function (req, status, err) {
				//console.log(req.responseText)
				console.log(req)

				console.log('Something went wrong', status, err);
				console.log(err)

			}
		});


	});

});