$(document).ready(function () {
	//$(".ajaxloader").fadeOut()
	var currentUrl = window.location.href;
	//getreport



	$(".refreshstatistic").each(function () {

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
                            SumSearch:"Количество кликов на поиске",
                            SessionDepthContext:"Глубина просмотра. РСЯ",
                            SessionDepthSearch:"Стоимость клика на поиске",
                            ShowsContext:"Количество показов в РСЯ",
                            ShowsSearch:"Глубина просмотра. Поиск",
                            StatDate:"Дата",
                            SumContext:"Количество показов на поиске",
                            GoalConversionSearch:"% целевых визитов. Поиск",
                            GoalConversionContext:"Стоимость клика в РСЯ",
							GoalCostSearch:"Стоимость ЛИДа. Поиск",
                            GoalCostContext:"% целевых визитов. РСЯ"
                        }
// 						"<th>ID РК</th>" +
//						"<th>Количество кликов в РСЯ</th>" +
//						"<th>Количество кликов на поиске</th>" +
//						"<th>Стоимость клика на поиске</th>" +
//						"<th>Глубина просмотра. РСЯ</th>" +
//						"<th>Глубина просмотра. Поиск</th>" +
//						"<th>Количество показов в РСЯ</th>" +
//						"<th>Количество показов на поиске</th>" +
//						"<th>Дата</th>" +
//						"<th>Стоимость клика в РСЯ</th>" +
//						"<th>% целевых визитов. Поиск</th>" +
//						"<th>% целевых визитов. РСЯ</th>" +
//						"<th>Стоимость ЛИДа. Поиск</th>" +
//						"<th>Стоимость ЛИДа. РСЯ</th>" +
//						"</tr></thead><tbody>"
//                         "<td>" + objsd.CampaignID + "</td>" +
//								"<td>" + objsd.ClicksContext + "</td>" +
//								"<td>" + objsd.ClicksSearch + "</td>" +
//								"<td>" + objsd.SumSearch + "</td>" +
//								"<td>" + objsd.SessionDepthContext + "</td>" +
//								"<td>" + objsd.SessionDepthSearch + "</td>" +
//								"<td>" + objsd.ShowsContext + "</td>" +
//								"<td>" + objsd.ShowsSearch + "</td>" +
//								"<td>" + objsd.StatDate + "</td>" +
//								"<td>" + objsd.SumContext + "</td>" +
//								"<td>" + objsd.GoalConversionSearch + "</td>" +
//								"<td>" + objsd.GoalConversionContext + "</p>" +
//								"<td>" + objsd.GoalCostSearch + "</td>" +
//								"<td>" + objsd.GoalCostContext + "</td>" + "</tr>")
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

//	$(".refreshstatistic").each(function RerfreshStat() {
//
//		$(this).click(function () {
//			$("#yandextable").empty();
//			$("#dateerror").empty();
//			$(".refreshstatistic").prop("disabled", true);
//			appendid = $(this).attr("result")
//			var startdate = $("#cal_1").val()
//			var enddate = $("#cal_2").val()
//			if (startdate === "" || enddate == "") {
//				console.log("Enter values")
//				$("#dateerror").append("<p style='border-radius: 2px; text-align: center; padding: 2px; border: solid 1px; color:red'>ВВедите даты</p>")
//				$(".refreshstatistic").prop('disabled', false);
//				return
//			}
//			console.log(startdate)
//			console.log(enddate)
//			//$(".refreshstatistic").hide(200)
//
//
//
//			$.ajax({
//				data: {
//					"startdate": startdate,
//					"enddate": enddate,
//				},
//				dataType: "json",
//				type: "POST",
//				url: currentUrl.replace("report", "getreport"),
//				beforeSend: function () {
//					$(".ajaxloader").fadeIn().show();
//					// $(".ajaxloader").html("<img id='imgcode' src='../static/img/loader.gif'/>");
//				},
//
//				success: function (data) {
//					$(".refreshstatistic").prop("disabled", false)
//					console.log("Main data recieved:", data)
//					if (data == null) {
//						console.log("Recieved data IS NULL")
//						$(".ajaxloader").fadeOut(0);
//						$("#yandextable").append("<p style='text-align: center; padding: 2px; font-size: 24px;'>В Яндексе нет данных за эти даты</p>");
//
//						return
//					}
//					console.log('Type of that data: ', typeof data)
//					var size = Object.keys(data).length
//					console.log('Object.keys(data).length: ', Object.keys(data).length)
//					var htmltab = "<table id='yandextableheaders' class='tablesorter' style='border-radius: 2px; text-align: center; padding: 2px; border: solid 1px'><thead><tr>" +
//						"<th>ID РК</th>" +
//						"<th>Количество кликов в РСЯ</th>" +
//						"<th>Количество кликов на поиске</th>" +
//						"<th>Стоимость клика на поиске</th>" +
//						"<th>Глубина просмотра. РСЯ</th>" +
//						"<th>Глубина просмотра. Поиск</th>" +
//						"<th>Количество показов в РСЯ</th>" +
//						"<th>Количество показов на поиске</th>" +
//						"<th>Дата</th>" +
//						"<th>Стоимость клика в РСЯ</th>" +
//						"<th>% целевых визитов. Поиск</th>" +
//						"<th>% целевых визитов. РСЯ</th>" +
//						"<th>Стоимость ЛИДа. Поиск</th>" +
//						"<th>Стоимость ЛИДа. РСЯ</th>" +
//						"</tr></thead><tbody>"
//					
//					$("#yandextable").append(htmltab)
//					for (accindex in data) {
//						//console.log('Recieved data type: ', typeof data)
//						//console.log("First loop (accindex in data) accindex = ", accindex)
//						//console.log('Type of that accindex: ', typeof accindex)
//						var size = Object.keys(accindex).length
//						
//							var objsd = data[accindex]
//							//var objsd = objsd3.Object
//							//console.log("Account object Data :", objsd)
//							//var objsd = data.Object
//							//console.log("The object: ", accData)
//							//console.log(typeof objsd)
//							//console.log(data[0].Data[0])
//							for (field in objsd) {
//								//console.log("Objsd: ",objsd)
//								//console.log("Field in objsd", field)
//								//console.log("value of objsd.field in objsd", objsd[field])
//								//console.log("type of objsd.field in objsd", typeof objsd[field])
//								//if (objsd[field] == "null" || "undefined" === typeof objsd.field) {
//								if (null === objsd[field]) {
//									objsd[field] = 0.00;
//									//console.log("objsd.field inside loop", objsd[field])
//								} else {
//									//console.log("Not null ", objsd[field])
//								};
//							}
//							//console.log("AFTER loop Objsd: ",objsd)
//							$("#yandextableheaders").append("<tr>" +
//								"<td>" + objsd.CampaignID + "</td>" +
//								"<td>" + objsd.ClicksContext + "</td>" +
//								"<td>" + objsd.ClicksSearch + "</td>" +
//								"<td>" + objsd.SumSearch + "</td>" +
//								"<td>" + objsd.SessionDepthContext + "</td>" +
//								"<td>" + objsd.SessionDepthSearch + "</td>" +
//								"<td>" + objsd.ShowsContext + "</td>" +
//								"<td>" + objsd.ShowsSearch + "</td>" +
//								"<td>" + objsd.StatDate + "</td>" +
//								"<td>" + objsd.SumContext + "</td>" +
//								"<td>" + objsd.GoalConversionSearch + "</td>" +
//								"<td>" + objsd.GoalConversionContext + "</p>" +
//								"<td>" + objsd.GoalCostSearch + "</td>" +
//								"<td>" + objsd.GoalCostContext + "</td>" + "</tr>")
//
//						}
//					
//					$("#yandextableheaders").append("</tbody></table>")
//					//$(".refreshstatistic").fadeIn(200)
//					{
//						$("#yandextableheaders").tablesorter();
//
//					}
//					$(".ajaxloader").fadeOut();
//				},
//				error: function (req, status, err) {
//					//console.log(req.responseText)
//					console.log(req)
//					console.log("json ", req.responseText)
//
//					console.log('Something went wrong', status, err);
//					console.log(err)
//					$(".refreshstatistic").prop('disabled', false);
//					$(".refreshstatistic").fadeIn(200)
//					$(".ajaxloader").fadeOut();
//					$("#yandextable").append(status, err)
//				}
//			});
//		});
//	})
