package data

type University struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var Universities = []University{
	{ID: "173499", Name: "Abdullah Gül Üniversitesi"},
	{ID: "326654", Name: "Acıbadem Mehmet Ali Aydınlar Üniversitesi"},
	{ID: "385368", Name: "Adana Alparslan Türkeş Bilim ve Teknoloji Üniversitesi"},
	{ID: "100259", Name: "Adıyaman Üniversitesi"},
	{ID: "100869", Name: "Afyon Kocatepe Üniversitesi"},
	{ID: "339999", Name: "Afyonkarahisar Sağlık Bilimleri Üniversitesi"},
	{ID: "101313", Name: "Ağrı İbrahim Çeçen Üniversitesi"},
	{ID: "101458", Name: "Akdeniz Üniversitesi"},
	{ID: "102005", Name: "Aksaray Üniversitesi"},
	{ID: "274883", Name: "Alanya Alaaddin Keykubat Üniversitesi"},
	{ID: "447835", Name: "Alanya Üniversitesi"},
	{ID: "326019", Name: "Altınbaş Üniversitesi"},
	{ID: "102198", Name: "Amasya Üniversitesi"},
	{ID: "102261", Name: "Anadolu Üniversitesi"},
	{ID: "411762", Name: "Ankara Bilim Üniversitesi"},
	{ID: "339985", Name: "Ankara Hacı Bayram Veli Üniversitesi"},
	{ID: "340001", Name: "Ankara Medipol Üniversitesi"},
	{ID: "324993", Name: "Ankara Müzik ve Güzel Sanatlar Üniversitesi"},
	{ID: "232846", Name: "Ankara Sosyal Bilimler Üniversitesi"},
	{ID: "102738", Name: "Ankara Üniversitesi"},
	{ID: "311737", Name: "Ankara Yıldırım Beyazıt Üniversitesi"},
	{ID: "447662", Name: "Antalya Belek Üniversitesi"},
	{ID: "327022", Name: "Antalya Bilim Üniversitesi"},
	{ID: "103400", Name: "Ardahan Üniversitesi"},
	{ID: "103443", Name: "Artvin Çoruh Üniversitesi"},
	{ID: "255896", Name: "Ataşehir Adıgüzel Meslek Yüksekokulu"},
	{ID: "103545", Name: "Atatürk Üniversitesi"},
	{ID: "104090", Name: "Atılım Üniversitesi"},
	{ID: "205351", Name: "Avrasya Üniversitesi"},
	{ID: "364158", Name: "Aydın Adnan Menderes Üniversitesi"},
	{ID: "104140", Name: "Bahçeşehir Üniversitesi"},
	{ID: "104213", Name: "Balıkesir Üniversitesi"},
	{ID: "274881", Name: "Bandırma Onyedi Eylül Üniversitesi"},
	{ID: "104578", Name: "Bartın Üniversitesi"},
	{ID: "104628", Name: "Başkent Üniversitesi"},
	{ID: "104803", Name: "Batman Üniversitesi"},
	{ID: "104922", Name: "Bayburt Üniversitesi"},
	{ID: "306383", Name: "Beykoz Üniversitesi"},
	{ID: "163894", Name: "Bezm-i Âlem Vakıf Üniversitesi"},
	{ID: "105024", Name: "Bilecik Şeyh Edebali Üniversitesi"},
	{ID: "105196", Name: "Bingöl Üniversitesi"},
	{ID: "251241", Name: "Biruni Üniversitesi"},
	{ID: "105248", Name: "Bitlis Eren Üniversitesi"},
	{ID: "105322", Name: "Boğaziçi Üniversitesi"},
	{ID: "341799", Name: "Bolu Abant İzzet Baysal Üniversitesi"},
	{ID: "355074", Name: "Burdur Mehmet Akif Ersoy Üniversitesi"},
	{ID: "173494", Name: "Bursa Teknik Üniversitesi"},
	{ID: "362386", Name: "Bursa Uludağ Üniversitesi"},
	{ID: "106499", Name: "Çağ Üniversitesi"},
	{ID: "106545", Name: "Çanakkale Onsekiz Mart Üniversitesi"},
	{ID: "107012", Name: "Çankaya Üniversitesi"},
	{ID: "107056", Name: "Çankırı Karatekin Üniversitesi"},
	{ID: "107211", Name: "Çukurova Üniversitesi"},
	{ID: "384591", Name: "Demiroğlu Bilim Üniversitesi"},
	{ID: "107723", Name: "Dicle Üniversitesi"},
	{ID: "108112", Name: "Doğuş Üniversitesi"},
	{ID: "108163", Name: "Dokuz Eylül Üniversitesi"},
	{ID: "109110", Name: "Düzce Üniversitesi"},
	{ID: "109290", Name: "Ege Üniversitesi"},
	{ID: "109868", Name: "Erciyes Üniversitesi"},
	{ID: "370189", Name: "Erzincan Binali Yıldırım Üniversitesi"},
	{ID: "173495", Name: "Erzurum Teknik Üniversitesi"},
	{ID: "110538", Name: "Eskişehir Osmangazi Üniversitesi"},
	{ID: "339997", Name: "Eskişehir Teknik Üniversitesi"},
	{ID: "163897", Name: "Fatih Sultan Mehmet Vakıf Üniversitesi"},
	{ID: "310017", Name: "Fenerbahçe Üniversitesi"},
	{ID: "110987", Name: "Fırat Üniversitesi"},
	{ID: "111395", Name: "Galatasaray Üniversitesi"},
	{ID: "133520", Name: "Gazi Üniversitesi"},
	{ID: "384577", Name: "Gaziantep İslam Bilim ve Teknoloji Üniversitesi"},
	{ID: "112080", Name: "Gaziantep Üniversitesi"},
	{ID: "260621", Name: "Gebze Teknik Üniversitesi"},
	{ID: "112836", Name: "Giresun Üniversitesi"},
	{ID: "113023", Name: "Gümüşhane Üniversitesi"},
	{ID: "113082", Name: "Hacettepe Üniversitesi"},
	{ID: "113681", Name: "Hakkari Üniversitesi"},
	{ID: "113699", Name: "Haliç Üniversitesi"},
	{ID: "113746", Name: "Harran Üniversitesi"},
	{ID: "138586", Name: "Hasan Kalyoncu Üniversitesi"},
	{ID: "367208", Name: "Hatay Mustafa Kemal Üniversitesi"},
	{ID: "114218", Name: "Hitit Üniversitesi"},
	{ID: "114315", Name: "Iğdır Üniversitesi"},
	{ID: "339998", Name: "Isparta Uygulamalı Bilimler Üniversitesi"},
	{ID: "114385", Name: "Işık Üniversitesi"},
	{ID: "274887", Name: "İbn Haldun Üniversitesi"},
	{ID: "105118", Name: "İhsan Doğramacı Bilkent Üniversitesi"},
	{ID: "114436", Name: "İnönü Üniversitesi"},
	{ID: "274882", Name: "İskenderun Teknik Üniversitesi"},
	{ID: "163900", Name: "İstanbul 29 Mayıs Üniversitesi"},
	{ID: "114773", Name: "İstanbul Arel Üniversitesi"},
	{ID: "339995", Name: "İstanbul Atlas Üniversitesi"},
	{ID: "114827", Name: "İstanbul Aydın Üniversitesi"},
	{ID: "448766", Name: "İstanbul Beykent Üniversitesi"},
	{ID: "114907", Name: "İstanbul Bilgi Üniversitesi"},
	{ID: "241174", Name: "İstanbul Esenyurt Üniversitesi"},
	{ID: "391144", Name: "İstanbul Galata Üniversitesi"},
	{ID: "315098", Name: "İstanbul Gedik Üniversitesi"},
	{ID: "130959", Name: "İstanbul Gelişim Üniversitesi"},
	{ID: "302687", Name: "İstanbul Kent Üniversitesi"},
	{ID: "115022", Name: "İstanbul Kültür Üniversitesi"},
	{ID: "173496", Name: "İstanbul Medeniyet Üniversitesi"},
	{ID: "163888", Name: "İstanbul Medipol Üniversitesi"},
	{ID: "447904", Name: "İstanbul Nişantaşı Üniversitesi"},
	{ID: "360777", Name: "İstanbul Okan Üniversitesi"},
	{ID: "274886", Name: "İstanbul Rumeli Üniversitesi"},
	{ID: "163898", Name: "İstanbul Sabahattin Zaim Üniversitesi"},
	{ID: "432690", Name: "İstanbul Sağlık ve Sosyal Bilimler Meslek Yüksekokulu"},
	{ID: "410560", Name: "İstanbul Sağlık ve Teknoloji Üniversitesi"},
	{ID: "220121", Name: "İstanbul Şişli Meslek Yüksekokulu"},
	{ID: "115069", Name: "İstanbul Teknik Üniversitesi"},
	{ID: "115335", Name: "İstanbul Ticaret Üniversitesi"},
	{ID: "440213", Name: "İstanbul Topkapı Üniversitesi"},
	{ID: "115373", Name: "İstanbul Üniversitesi"},
	{ID: "339984", Name: "İstanbul Üniversitesi-Cerrahpaşa"},
	{ID: "315415", Name: "İstanbul Yeni Yüzyıl Üniversitesi"},
	{ID: "274888", Name: "İstinye Üniversitesi"},
	{ID: "302686", Name: "İzmir Bakırçay Üniversitesi"},
	{ID: "302685", Name: "İzmir Demokrasi Üniversitesi"},
	{ID: "116147", Name: "İzmir Ekonomi Üniversitesi"},
	{ID: "173498", Name: "İzmir Katip Çelebi Üniversitesi"},
	{ID: "334490", Name: "İzmir Kavram Meslek Yüksekokulu"},
	{ID: "339996", Name: "İzmir Tınaztepe Üniversitesi"},
	{ID: "116207", Name: "İzmir Yüksek Teknoloji Enstitüsü"},
	{ID: "116281", Name: "Kadir Has Üniversitesi"},
	{ID: "116345", Name: "Kafkas Üniversitesi"},
	{ID: "339994", Name: "Kahramanmaraş İstiklal Üniversitesi"},
	{ID: "116608", Name: "Kahramanmaraş Sütçü İmam Üniversitesi"},
	{ID: "325756", Name: "Kapadokya Üniversitesi"},
	{ID: "116950", Name: "Karabük Üniversitesi"},
	{ID: "117127", Name: "Karadeniz Teknik Üniversitesi"},
	{ID: "117553", Name: "Karamanoğlu Mehmetbey Üniversitesi"},
	{ID: "117673", Name: "Kastamonu Üniversitesi"},
	{ID: "339993", Name: "Kayseri Üniversitesi"},
	{ID: "117803", Name: "Kırıkkale Üniversitesi"},
	{ID: "118122", Name: "Kırklareli Üniversitesi"},
	{ID: "354265", Name: "Kırşehir Ahi Evran Üniversitesi"},
	{ID: "118186", Name: "Kilis 7 Aralık Üniversitesi"},
	{ID: "411763", Name: "Kocaeli Sağlık ve Teknoloji Üniversitesi"},
	{ID: "118239", Name: "Kocaeli Üniversitesi"},
	{ID: "118853", Name: "Koç Üniversitesi"},
	{ID: "241176", Name: "Konya Gıda ve Tarım Üniversitesi"},
	{ID: "339979", Name: "Konya Teknik Üniversitesi"},
	{ID: "166433", Name: "KTO Karatay Üniversitesi"},
	{ID: "351149", Name: "Kütahya Dumlupınar Üniversitesi"},
	{ID: "339982", Name: "Kütahya Sağlık Bilimleri Üniversitesi"},
	{ID: "332474", Name: "Lokman Hekim Üniversitesi"},
	{ID: "339983", Name: "Malatya Turgut Özal Üniversitesi"},
	{ID: "118883", Name: "Maltepe Üniversitesi"},
	{ID: "315839", Name: "Manisa Celâl Bayar Üniversitesi"},
	{ID: "118994", Name: "Mardin Artuklu Üniversitesi"},
	{ID: "119094", Name: "Marmara Üniversitesi"},
	{ID: "215913", Name: "Mef Üniversitesi"},
	{ID: "119917", Name: "Mersin Üniversitesi"},
	{ID: "120301", Name: "Mimar Sinan Güzel Sanatlar Üniversitesi"},
	{ID: "442563", Name: "Mudanya Üniversitesi"},
	{ID: "120444", Name: "Muğla Sıtkı Koçman Üniversitesi"},
	{ID: "307919", Name: "Munzur Üniversitesi"},
	{ID: "121164", Name: "Muş Alparslan Üniversitesi"},
	{ID: "173500", Name: "Necmettin Erbakan Üniversitesi"},
	{ID: "246224", Name: "Nevşehir Hacı Bektaş Veli Üniversitesi"},
	{ID: "306556", Name: "Niğde Ömer Halisdemir Üniversitesi"},
	{ID: "163891", Name: "Nuh Naci Yazgan Üniversitesi"},
	{ID: "121946", Name: "Ondokuz Mayıs Üniversitesi"},
	{ID: "122395", Name: "Ordu Üniversitesi"},
	{ID: "122571", Name: "Orta Doğu Teknik Üniversitesi"},
	{ID: "122735", Name: "Osmaniye Korkut Ata Üniversitesi"},
	{ID: "324992", Name: "Ostim Teknik Üniversitesi"},
	{ID: "122827", Name: "Özyeğin Üniversitesi"},
	{ID: "122831", Name: "Pamukkale Üniversitesi"},
	{ID: "136233", Name: "Piri Reis Üniversitesi"},
	{ID: "123221", Name: "Recep Tayyip Erdoğan Üniversitesi"},
	{ID: "123400", Name: "Sabancı Üniversitesi"},
	{ID: "270121", Name: "Sağlık Bilimleri Üniversitesi"},
	{ID: "339988", Name: "Sakarya Uygulamalı Bilimler Üniversitesi"},
	{ID: "123409", Name: "Sakarya Üniversitesi"},
	{ID: "339989", Name: "Samsun Üniversitesi"},
	{ID: "241177", Name: "Sanko Üniversitesi"},
	{ID: "123902", Name: "Selçuk Üniversitesi"},
	{ID: "124703", Name: "Siirt Üniversitesi"},
	{ID: "124805", Name: "Sinop Üniversitesi"},
	{ID: "339990", Name: "Sivas Bilim ve Teknoloji Üniversitesi"},
	{ID: "344737", Name: "Sivas Cumhuriyet Üniversitesi"},
	{ID: "124902", Name: "Süleyman Demirel Üniversitesi"},
	{ID: "125536", Name: "Şırnak Üniversitesi"},
	{ID: "339991", Name: "Tarsus Üniversitesi"},
	{ID: "163892", Name: "Ted Üniversitesi"},
	{ID: "356278", Name: "Tekirdağ Namık Kemal Üniversitesi"},
	{ID: "125552", Name: "TOBB Ekonomi ve Teknoloji Üniversitesi"},
	{ID: "367245", Name: "Tokat Gaziosmanpaşa Üniversitesi"},
	{ID: "163889", Name: "Toros Üniversitesi"},
	{ID: "339992", Name: "Trabzon Üniversitesi"},
	{ID: "125577", Name: "Trakya Üniversitesi"},
	{ID: "203267", Name: "Türk Hava Kurumu Üniversitesi"},
	{ID: "201784", Name: "Türk-Alman Üniversitesi"},
	{ID: "125968", Name: "Ufuk Üniversitesi"},
	{ID: "126537", Name: "Uşak Üniversitesi"},
	{ID: "206795", Name: "Üsküdar Üniversitesi"},
	{ID: "337414", Name: "Van Yüzüncü Yıl Üniversitesi"},
	{ID: "126742", Name: "Yalova Üniversitesi"},
	{ID: "126773", Name: "Yaşar Üniversitesi"},
	{ID: "126818", Name: "Yeditepe Üniversitesi"},
	{ID: "126982", Name: "Yıldız Teknik Üniversitesi"},
	{ID: "359730", Name: "Yozgat Bozok Üniversitesi"},
	{ID: "206792", Name: "Yüksek İhtisas Üniversitesi"},
	{ID: "365890", Name: "Zonguldak Bülent Ecevit Üniversitesi"},
}

func IsValidUniversityID(id string) bool {
	for _, univ := range Universities {
		if univ.ID == id {
			return true
		}
	}
	return false
}
