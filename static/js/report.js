$(document).ready(function () {
	//$(".ajaxloader").fadeOut()
	var currentUrl = window.location.href;
	//getreport



	$(".refreshstatistic").each(function () {
		// if in currentUrl contains getaccountstat, than use logic for one account
		if (currentUrl.includes("getaccountstat")) {
            console.log("getaccountstat logic used")
            $(this).click(function () {
                $("#yandextable").empty();
                $("#dateerror").empty();
                $(".refreshstatistic").prop("disabled", true);
                appendid = $(this).attr("result")
                var startdate = $("#cal_1").val()
                var enddate = $("#cal_2").val()
                if (startdate === "" || enddate == "") {
                    console.log("Enter values")
                    $("#dateerror").append("<p style='border-radius: 2px; text-align: center; padding: 2px; border: solid 1px; color:red'>Введите даты</p>")
                    $(".refreshstatistic").prop('disabled', false);
                    return
                }
                console.log(startdate)
                console.log(enddate)
                //$(".refreshstatistic").hide(200)
                $.ajax({
                    data: {
                        "startdate": startdate,
                        "enddate": enddate,
                    },
                    dataType: "json",
                    type: "POST",
                    url: currentUrl,
                    beforeSend: function () {
                        $(".ajaxloader").fadeIn().show();
                        // $(".ajaxloader").html("<img id='imgcode' src='../static/img/loader.gif'/>");
                    },
                    success: function (data) {
                        $(".refreshstatistic").prop("disabled", false)
                        console.log("Main data recieved:", data)
						//data =
                        if (data == null) {
                            console.log("Recieved data IS NULL")
                            $(".ajaxloader").fadeOut(0);
                            $("#yandextable").append("<p style='text-align: center; padding: 2px; font-size: 24px;'>В Яндексе нет данных за эти даты</p>");
                            return
                        }
                        // for (accindex in data) {
                        //     var objsd = data[accindex]
                        //     for (field in objsd) {
                        //         if (null === objsd[field] || objsd[field] == undefined || objsd[field] == "") {
                        //             objsd[field] = 0;
                        //         }
                        //     }
                        // }
                        console.log('Type of that data: ', typeof data);
                        //data = data.response;
                        // var m = []
                        // for (i in data) {
                        // 	var objs2 = data[i];
                        //     for (y in objs2.stats) {
                        //         m.push(objs2.stats[y])
							// }
                        // }
                        // console.log("Debug m ",m);
                        var table = new KingTable({
                            element: document.getElementById("yandextable"),
                            data: data,

                        });
                        table.render();

                        console.log('All data : ', data)
                        $(".ajaxloader").fadeOut();
                    },
                    error: function (req, status, err) {
                        //console.log(req.responseText)
                        console.log(req);
                        console.log("json ", req.responseText);
                        console.log('Something went wrong', status, err);
                        console.log(err);
                        $(".refreshstatistic").prop('disabled', false);
                        $(".refreshstatistic").fadeIn(200)
                        $(".ajaxloader").fadeOut();
                        $("#yandextable").append(status, err)
                    }
                });
            });
            return
        }
		$(this).click(function () {
			$("#yandextable").empty();
			$("#dateerror").empty();
			$(".refreshstatistic").prop("disabled", true);
			appendid = $(this).attr("result")
			var startdate = $("#cal_1").val()
			var enddate = $("#cal_2").val()
			if (startdate === "" || enddate == "") {
				console.log("Enter values")
				$("#dateerror").append("<p style='border-radius: 2px; text-align: center; padding: 2px; border: solid 1px; color:red'>ВВедите даты</p>")
				$(".refreshstatistic").prop('disabled', false);
				return
			}
			console.log(startdate)
			console.log(enddate)
			//$(".refreshstatistic").hide(200)
			$.ajax({
				data: {
					"startdate": startdate,
					"enddate": enddate,
				},
				dataType: "json",
				type: "POST",
				url: currentUrl.replace("report", "getreport"),
				beforeSend: function () {
					$(".ajaxloader").fadeIn().show();
					// $(".ajaxloader").html("<img id='imgcode' src='../static/img/loader.gif'/>");
				},

				success: function (data) {
					$(".refreshstatistic").prop("disabled", false)
					console.log("Main data recieved:", data)
					if (data == null) {
						console.log("Recieved data IS NULL")
						$(".ajaxloader").fadeOut(0);
						$("#yandextable").append("<p style='text-align: center; padding: 2px; font-size: 24px;'>В Яндексе нет данных за эти даты</p>");
						return
					}
					for (accindex in data) {
						var objsd = data[accindex]
						for (field in objsd) {
							if (null === objsd[field] || objsd[field] == undefined || objsd[field] == "") {
								objsd[field] = 0;
							}
						}
					}
					console.log('Type of that data: ', typeof data);
					var size = Object.keys(data).length;
					var table = new KingTable({
						element: document.getElementById("yandextable"),
						data: data,
                        columns: {
                            CampaignID: "ID РК",
                            ClicksContext:"Количество кликов в РСЯ",
                            ClicksSearch:"Количество кликов на поиске",
                            SumSearch:"Стоимость кликов на поиске.",
                            SessionDepthContext:"Глубина просмотра. РСЯ",
                            SessionDepthSearch:"Глубина просмотра. Поиск",
                            ShowsContext:"Количество показов в РСЯ",
                            ShowsSearch:"Количество показов в Поиск",
                            StatDate:"Дата",
                            SumContext:"Стоимость кликов в РСЯ.",
                            GoalConversionSearch:"% целевых визитов. Поиск",
                            GoalConversionContext:"% целевых визитов. РСЯ",
							GoalCostSearch:"Стоимость ЛИДа. Поиск",
                            GoalCostContext:"Стоимость ЛИДа. РСЯ"
                        }
					});
					table.render();
					console.log('All data : ', data)
					$(".ajaxloader").fadeOut();
				},
				error: function (req, status, err) {
					//console.log(req.responseText)
					console.log(req);
					console.log("json ", req.responseText);
					console.log('Something went wrong', status, err);
					console.log(err);
					$(".refreshstatistic").prop('disabled', false);
					$(".refreshstatistic").fadeIn(200)
					$(".ajaxloader").fadeOut();
					$("#yandextable").append(status, err)
				}
			});
		});
	})
});
