package yandex

//
//func collectCampaingsfromAgency(login, token string) (res []campaigns.GetResponse, err error) {
//	fmt.Println("collectCampaingsfromAgency START", time.Now())
//	resultA, err := collectAgencyClients(login, token)
//	if err != nil {
//		logrus.Errorln("collectCampaingsfromAgency  error: ", err)
//		return
//	}
//	var camps []campaigns.GetResponse
//	for _, ag := range resultA.Clients {
//		for _, login := range ag.Representatives {
//			result, err := collectCampaings(login.Login, token)
//			if err != nil {
//				logrus.Errorln("cant collect campaings with parameters: collectCampaings(%s, %s) error: %v", login.Login, token, err)
//				return nil, fmt.Errorf("cant collect campaings with parameters: collectCampaings(%s, %s), error: %v", login.Login, token, err)
//			}
//			camps = append(camps, result)
//		}
//	}
//	fmt.Println("collectCampaingsfromAgency END", time.Now())
//	return camps, nil
//}
//func collectCampaingsfromAgencyConcurently(login, token string) (res []campaigns.GetResponse, err error) {
//	fmt.Println("collectCampaingsfromAgencyConcurently START", time.Now())
//
//	var YandexConnectionsLimit = 5
//	chAC := make(chan gc.ClientGetItem, 4) // This number 3 can be anything as long as it's larger than YandexConnectionsLimit
//	var wg sync.WaitGroup
//	resultA, err := collectAgencyClients(login, token)
//	if err != nil {
//		logrus.Errorln("collectCampaingsfromAgency  error: ", err)
//		return
//	}
//	logrus.Debug("Debug 2 ")
//	var camps []campaigns.GetResponse
//	logrus.Debug("Debug 3 ")
//	for i := 0; i < YandexConnectionsLimit; i++ {
//		logrus.Info("Iterating through YandexConnectionsLimit")
//		wg.Add(1)
//		go func() {
//			for {
//				login, ok := <-chAC
//				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
//					wg.Done()
//					return
//				}
//				logrus.Info("Inside gorouitne with login:", login)
//				result, err := collectCampaings(login.Login, token)
//				if err != nil {
//					logrus.Errorln("cant collect campaings with parameters: collectCampaings(%s, %s) error: %v", login, token, err)
//					//return nil, fmt.Errorf("cant collect campaings with parameters: collectCampaings(%s, %s), error: %v", login.Login, token, err)
//				}
//				camps = append(camps, result)
//				logrus.Info("Inside gorouitne append is ok for login:", login)
//			}
//		}()
//	}
//	for _, c := range resultA.Clients {
//
//		chAC <- c // add agClient to the queue
//
//	}
//	close(chAC) // This tells the goroutines there's nothing else to do
//	wg.Wait()   // Wait for the threads to finish
//
//	fmt.Println("collectCampaingsfromAgency END", time.Now())
//	return camps, nil
//}
//
//func collectCampaingsfromAgencyConcurently2(login, token string) (res []campaigns.GetResponse, err error) {
//	fmt.Println("collectCampaingsfromAgencyConcurently START", time.Now())
//	var YandexConnectionsLimit = 4
//	var wg sync.WaitGroup
//	wg.Add(YandexConnectionsLimit)
//	resultA, err := collectAgencyClients(login, token)
//	if err != nil {
//		logrus.Errorln("collectCampaingsfromAgency  error: ", err)
//		return
//	}
//	var camps []campaigns.GetResponse
//	for _, ag := range resultA.Clients {
//		go func() {
//			for {
//				for _, login := range ag.Representatives {
//					result, err := collectCampaings(login.Login, token)
//					if err != nil {
//						logrus.Errorln("cant collect campaings with parameters: collectCampaings(%s, %s) error: %v", login.Login, token, err)
//						//return nil, fmt.Errorf("cant collect campaings with parameters: collectCampaings(%s, %s), error: %v", login.Login, token, err)
//					}
//					camps = append(camps, result)
//				}
//				wg.Done()
//			}
//		}()
//	}
//	wg.Wait() // Wait for the threads to finish
//
//	fmt.Println("collectCampaingsfromAgency END", time.Now())
//	return camps, nil
//}
