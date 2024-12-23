package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	host                     = os.Getenv("HOST")
	port                     = os.Getenv("PORT")
	user                     = os.Getenv("USER")
	password                 = os.Getenv("PASSWORD")
	dbname                   = os.Getenv("DBNAME")
	store_count              = 101
	order_count              = 0
	initial_order_count_k, _ = strconv.Atoi(os.Getenv("INITIAL_ORDER_COUNT_K"))
	order_update_count_k, _  = strconv.Atoi(os.Getenv("ORDER_UPDATE_COUNT_K"))
	new_order_count_k, _     = strconv.Atoi(os.Getenv("NEW_ORDER_COUNT_K"))

	order_statuses = []string{"Created", "Processing", "Shipped", "Delivered", "Cancelled"}
	stores         = []StoreData{}
)

type OrderData struct {
	OrderAmount   float32 `json:"order_amount"`
	OrderStatusId int     `json:"order_status_id"`
	OrderStatus   string  `json:"order_status"`
}

type Loc struct {
	Lat  float32 `json:"store_lat"`
	Long float32 `json:"store_long"`
}

type StoreData struct {
	StoreId   int    `json:"store_id"`
	StoreAddr string `json:"store_addr"`
	StoreLoc  Loc    `json:"store_loc"`
}

type Data struct {
	Order OrderData `json:"order"`
	Store StoreData `json:"store"`
}

func main() {
	fmt.Println("INITIAL_ORDER_COUNT_K: ", initial_order_count_k)
	fmt.Println("ORDER_UPDATE_COUNT_K: ", order_update_count_k)
	fmt.Println("NEW_ORDER_COUNT_K: ", new_order_count_k)

	start := time.Now()
	generateStoreArray()
	generateOrders(initial_order_count_k)
	fmt.Println("Finished generating initial orders in", time.Since(start).Seconds(), "seconds.")
	order_count = initial_order_count_k

	for {
		start = time.Now()
		updateOrders()
		fmt.Println("Finished updating orders in", time.Since(start).Seconds(), "seconds.")
		time.Sleep(250 * time.Millisecond)

		start = time.Now()
		generateOrders(new_order_count_k)
		fmt.Println("Finished generating new orders in", time.Since(start).Seconds(), "seconds.")
		order_count += new_order_count_k

	}
}

func generateOrders(order_count int) {
	psqlInfo := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbname, user, password)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(500)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// generate orders concurrently
	var wg sync.WaitGroup
	for i := 0; i < order_count; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {

				store := stores[rand.Intn(store_count)]

				order := OrderData{float32(rand.Intn(18001)+2000) / 100, 1, "Created"}
				data, _ := json.Marshal(Data{order, store})
				query := `INSERT INTO order_schema.order (data) VALUES ($1);`
				_, err = db.Exec(query, data)
				if err != nil {
					panic(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func updateOrders() {
	psqlInfo := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbname, user, password)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// get a random int between 100 and 300 to cancel orders
	cancel_order_count := rand.Intn(300-100) + 100

	// cancel orders
	_, err = db.Exec(`UPDATE order_schema.order
                      SET data = jsonb_set(jsonb_set("data", '{order, order_status_id}', to_jsonb(5)), '{order, order_status}', '"Cancelled"'), updated_at = CURRENT_TIMESTAMP
                      WHERE order_id IN (SELECT order_id
                                         FROM order_schema.order
                                         WHERE data->'order'->'order_status_id' NOT IN (to_jsonb(4),to_jsonb(5))
                                         ORDER BY RANDOM()
                                         LIMIT $1);`, cancel_order_count)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("orders have been cancelled")
	}
	// update remaining orders
	// get list of order_ids to update
	rows, err := db.Query(`SELECT order_id, data->'order'->'order_status_id'
                                FROM order_schema.order
                                WHERE data->'order'->'order_status_id' NOT IN (to_jsonb(4),to_jsonb(5))
                                ORDER BY RANDOM()
                                LIMIT $1;`, order_update_count_k*1000-cancel_order_count)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("order status data retrieved from db")
	}
	order_ids := []int{}
	statuses := []int{}
	var id int
	var status int
	for rows.Next() {
		rows.Scan(&id, &status)
		order_ids = append(order_ids, id)
		statuses = append(statuses, status+1)
	}

	for i := range order_ids {
		sql := fmt.Sprintf(`UPDATE order_schema.order
                      SET data = jsonb_set(jsonb_set("data", '{order, order_status_id}', to_jsonb(%v)), '{order, order_status}', '"%v"'),
                          updated_at = CURRENT_TIMESTAMP
                      WHERE order_id = %v;`, statuses[i], order_statuses[statuses[i]-1], order_ids[i])
		_, err = db.Exec(sql)
		if err != nil {
			panic(err)
		}
	}
}

func generateStoreArray() {

	stores = append(stores, StoreData{1, "1600 Pennsylvania Ave NW, Washington, DC, 20500", Loc{38.897957, -77.036560}})
	stores = append(stores, StoreData{2, "20 W 34th St, New York, NY, 10001", Loc{40.748440, -73.985664}})                // Empire State Building
	stores = append(stores, StoreData{3, "Liberty Island, New York, NY, 10004", Loc{40.689249, -74.044500}})              // Statue of Liberty
	stores = append(stores, StoreData{4, "46th St & Broadway, New York, NY, 10036", Loc{40.758000, -73.985500}})          // Times Square
	stores = append(stores, StoreData{5, "Central Park West & 72nd St, New York, NY, 10023", Loc{40.774825, -73.974052}}) // Central Park (Dakota area)
	stores = append(stores, StoreData{6, "45 Rockefeller Plaza, New York, NY, 10111", Loc{40.758700, -73.978700}})        // Rockefeller Center
	stores = append(stores, StoreData{7, "285 Fulton St, New York, NY, 10007", Loc{40.712742, -74.013382}})               // One World Trade Center
	stores = append(stores, StoreData{8, "89 E 42nd St, New York, NY, 10017", Loc{40.752726, -73.977229}})                // Grand Central Terminal
	stores = append(stores, StoreData{9, "200 Eastern Pkwy, Brooklyn, NY, 11238", Loc{40.671206, -73.963631}})            // Brooklyn Museum
	stores = append(stores, StoreData{10, "11 W 53rd St, New York, NY, 10019", Loc{40.761433, -73.977622}})               // MoMA

	stores = append(stores, StoreData{11, "1000 5th Ave, New York, NY, 10028", Loc{40.779437, -73.963244}})                  // The Met
	stores = append(stores, StoreData{12, "200 Central Park West, New York, NY, 10024", Loc{40.781324, -73.973988}})         // American Museum of Natural History
	stores = append(stores, StoreData{13, "911 Memorial Plaza, New York, NY, 10007", Loc{40.711476, -74.013310}})            // 9/11 Memorial & Museum
	stores = append(stores, StoreData{14, "High Line, New York, NY, 10011", Loc{40.747993, -74.004765}})                     // The High Line
	stores = append(stores, StoreData{15, "Hudson Yards, 20 Hudson Yards, New York, NY, 10001", Loc{40.753826, -74.002233}}) // The Vessel
	stores = append(stores, StoreData{16, "Golden Gate Bridge, San Francisco, CA, 94129", Loc{37.819928, -122.478255}})
	stores = append(stores, StoreData{17, "Alcatraz Island, San Francisco, CA, 94133", Loc{37.826977, -122.422956}})
	stores = append(stores, StoreData{18, "Pier 39, San Francisco, CA, 94133", Loc{37.808674, -122.409821}})         // Fisherman’s Wharf
	stores = append(stores, StoreData{19, "1070 Lombard St, San Francisco, CA, 94109", Loc{37.802139, -122.418740}}) // Lombard Street
	stores = append(stores, StoreData{20, "501 Stanyan St, San Francisco, CA, 94117", Loc{37.770202, -122.454478}})  // Golden Gate Park

	stores = append(stores, StoreData{21, "1313 Disneyland Dr, Anaheim, CA, 92802", Loc{33.812092, -117.918974}})              // Disneyland
	stores = append(stores, StoreData{22, "Hollywood Blvd & Vine St, Los Angeles, CA, 90028", Loc{34.101998, -118.326519}})    // Hollywood Walk of Fame
	stores = append(stores, StoreData{23, "2800 E Observatory Rd, Los Angeles, CA, 90027", Loc{34.118434, -118.300393}})       // Griffith Observatory
	stores = append(stores, StoreData{24, "100 Universal City Plaza, Universal City, CA, 91608", Loc{34.138116, -118.353378}}) // Universal Studios Hollywood
	stores = append(stores, StoreData{25, "200 Santa Monica Pier, Santa Monica, CA, 90401", Loc{34.009351, -118.497468}})      // Santa Monica Pier
	stores = append(stores, StoreData{26, "1200 Getty Center Dr, Los Angeles, CA, 90049", Loc{34.078036, -118.474095}})        // The Getty Center
	stores = append(stores, StoreData{27, "3400 W Riverside Dr, Burbank, CA, 91505", Loc{34.154320, -118.337668}})             // Warner Bros. Studio Tour
	stores = append(stores, StoreData{28, "Los Angeles, CA, 90068", Loc{34.134115, -118.321548}})                              // Hollywood Sign (approx.)
	stores = append(stores, StoreData{29, "1111 S Figueroa St, Los Angeles, CA, 90015", Loc{34.043017, -118.267254}})          // Crypto.com Arena
	stores = append(stores, StoreData{30, "111 S Grand Ave, Los Angeles, CA, 90012", Loc{34.055345, -118.249845}})             // Walt Disney Concert Hall

	stores = append(stores, StoreData{31, "Redwood National and State Parks, Orick, CA, 95555", Loc{41.213177, -124.004628}})
	stores = append(stores, StoreData{32, "9035 Village Dr, Yosemite Valley, CA, 95389", Loc{37.746935, -119.589380}})      // Yosemite Visitor Center
	stores = append(stores, StoreData{33, "6554 Park Blvd, Joshua Tree, CA, 92252", Loc{34.076733, -116.305974}})           // Joshua Tree NP
	stores = append(stores, StoreData{34, "Entrance near CA-190, Death Valley, CA, 92328", Loc{36.241711, -116.825833}})    // Death Valley NP
	stores = append(stores, StoreData{35, "Hoover Dam Access Rd, Boulder City, NV, 89005", Loc{36.016112, -114.737778}})    // Hoover Dam
	stores = append(stores, StoreData{36, "20 South Entrance Rd, Grand Canyon, AZ, 86023", Loc{36.052987, -112.121483}})    // Grand Canyon (South Rim)
	stores = append(stores, StoreData{37, "Navajo Tribal Park, Page, AZ, 86040", Loc{36.861897, -111.374438}})              // Antelope Canyon (approx.)
	stores = append(stores, StoreData{38, "US-163 Scenic, Oljato-Monument Valley, UT, 84536", Loc{37.004331, -110.164668}}) // Monument Valley
	stores = append(stores, StoreData{39, "Bryce Canyon National Park, UT-63, Bryce, UT, 84764", Loc{37.593038, -112.187089}})
	stores = append(stores, StoreData{40, "1 Zion Park Blvd, Springdale, UT, 84767", Loc{37.198667, -112.986146}}) // Zion NP

	stores = append(stores, StoreData{41, "Arches National Park, Moab, UT, 84532", Loc{38.733082, -109.592514}})
	stores = append(stores, StoreData{42, "1000 US Hwy 36, Estes Park, CO, 80517", Loc{40.365742, -105.556157}})              // Rocky Mtn NP (Beaver Meadows)
	stores = append(stores, StoreData{43, "13000 Hwy 244, Keystone, SD, 57751", Loc{43.879102, -103.459067}})                 // Mount Rushmore
	stores = append(stores, StoreData{44, "12151 Avenue of the Chiefs, Crazy Horse, SD, 57730", Loc{43.820507, -103.640316}}) // Crazy Horse Memorial
	stores = append(stores, StoreData{45, "25216 Ben Reifel Rd, Interior, SD, 57750", Loc{43.771193, -101.928856}})           // Badlands NP
	stores = append(stores, StoreData{46, "11 N 4th St, St. Louis, MO, 63102", Loc{38.625927, -90.189128}})                   // Gateway Arch
	stores = append(stores, StoreData{47, "60 E Broadway, Bloomington, MN, 55425", Loc{44.855540, -93.242400}})               // Mall of America
	stores = append(stores, StoreData{48, "233 S Wacker Dr, Chicago, IL, 60606", Loc{41.878876, -87.635915}})                 // Willis (Sears) Tower
	stores = append(stores, StoreData{49, "201 E Randolph St, Chicago, IL, 60602", Loc{41.885003, -87.622908}})               // Millennium Park
	stores = append(stores, StoreData{50, "332 Prospect St, Niagara Falls, NY, 14303", Loc{43.084358, -79.067834}})           // Niagara Falls (NY side)

	stores = append(stores, StoreData{51, "1 Cedar Point Dr, Sandusky, OH, 44870", Loc{41.478424, -82.678206}})              // Cedar Point
	stores = append(stores, StoreData{52, "20900 Oakwood Blvd, Dearborn, MI, 48124", Loc{42.300158, -83.233554}})            // The Henry Ford
	stores = append(stores, StoreData{53, "7274 Main St, Mackinac Island, MI, 49757", Loc{45.845569, -84.617379}})           // Mackinac Island
	stores = append(stores, StoreData{54, "1 Lodge St, Asheville, NC, 28803", Loc{35.540495, -82.553663}})                   // Biltmore Estate
	stores = append(stores, StoreData{55, "107 Park Headquarters Rd, Gatlinburg, TN, 37738", Loc{35.686777, -83.536853}})    // Great Smoky Mountains
	stores = append(stores, StoreData{56, "2700 Dollywood Parks Blvd, Pigeon Forge, TN, 37863", Loc{35.795065, -83.535606}}) // Dollywood
	stores = append(stores, StoreData{57, "300 Alamo Plaza, San Antonio, TX, 78205", Loc{29.425967, -98.486141}})            // The Alamo
	stores = append(stores, StoreData{58, "2101 E NASA Pkwy, Houston, TX, 77058", Loc{29.551832, -95.097686}})               // NASA Johnson Space Center
	stores = append(stores, StoreData{59, "849 E Commerce St, San Antonio, TX, 78205", Loc{29.424349, -98.484184}})          // San Antonio River Walk
	stores = append(stores, StoreData{60, "Bourbon St, New Orleans, LA, 70116", Loc{29.964111, -90.061628}})                 // Bourbon Street

	stores = append(stores, StoreData{61, "French Quarter, New Orleans, LA, 70116", Loc{29.958443, -90.065007}})
	stores = append(stores, StoreData{62, "3764 Elvis Presley Blvd, Memphis, TN, 38116", Loc{35.046251, -90.026891}}) // Graceland
	stores = append(stores, StoreData{63, "Beale St, Memphis, TN, 38103", Loc{35.139170, -90.050554}})                // Beale Street
	stores = append(stores, StoreData{64, "222 5th Ave S, Nashville, TN, 37203", Loc{36.159166, -86.776853}})         // Country Music Hall of Fame
	stores = append(stores, StoreData{65, "2804 Opryland Dr, Nashville, TN, 37214", Loc{36.206611, -86.693369}})      // Grand Ole Opry
	stores = append(stores, StoreData{66, "1 Mammoth Cave Pkwy, Mammoth Cave, KY, 42259", Loc{37.187042, -86.100843}})
	stores = append(stores, StoreData{67, "700 Central Ave, Louisville, KY, 40208", Loc{38.206300, -85.759532}})          // Churchill Downs
	stores = append(stores, StoreData{68, "1180 Seven Seas Dr, Lake Buena Vista, FL, 32830", Loc{28.417663, -81.581212}}) // Walt Disney World
	stores = append(stores, StoreData{69, "6000 Universal Blvd, Orlando, FL, 32819", Loc{28.474321, -81.467819}})         // Universal Orlando
	stores = append(stores, StoreData{70, "Space Commerce Way, Merritt Island, FL, 32899", Loc{28.524302, -80.650382}})   // Kennedy Space Center

	stores = append(stores, StoreData{71, "40001 State Road 9336, Homestead, FL, 33034", Loc{25.392534, -80.583095}})                  // Everglades NP
	stores = append(stores, StoreData{72, "Ocean Dr & 5th St, Miami Beach, FL, 33139", Loc{25.772414, -80.132279}})                    // South Beach (approx.)
	stores = append(stores, StoreData{73, "1801 W International Speedway Blvd, Daytona Beach, FL, 32114", Loc{29.188544, -81.070575}}) // Daytona Intl Speedway
	stores = append(stores, StoreData{74, "101 Visitor Center Dr, Williamsburg, VA, 23185", Loc{37.270587, -76.693722}})               // Colonial Williamsburg
	stores = append(stores, StoreData{75, "1 Memorial Ave, Fort Myer, VA, 22211", Loc{38.876663, -77.070874}})                         // Arlington Cemetery
	stores = append(stores, StoreData{76, "3200 Mount Vernon Memorial Hwy, Mount Vernon, VA, 22121", Loc{38.729466, -77.107667}})      // George Washington’s Mt Vernon
	stores = append(stores, StoreData{77, "931 Thomas Jefferson Pkwy, Charlottesville, VA, 22902", Loc{38.009537, -78.450354}})        // Monticello
	stores = append(stores, StoreData{78, "3655 US Hwy 211 E, Luray, VA, 22835", Loc{38.662086, -78.372244}})                          // Skyline Drive (Shenandoah NP)
	stores = append(stores, StoreData{79, "2834 Washington Rd, Gettysburg, PA, 17325", Loc{39.816299, -77.232071}})                    // Gettysburg NMP (VC)
	stores = append(stores, StoreData{80, "1 Audrey Zapp Dr, Jersey City, NJ, 07305", Loc{40.707151, -74.045511}})                     // Liberty State Park

	stores = append(stores, StoreData{81, "1 Borgata Way, Atlantic City, NJ, 08401", Loc{39.379496, -74.451290}})      // Borgata
	stores = append(stores, StoreData{82, "520 Chestnut St, Philadelphia, PA, 19106", Loc{39.949610, -75.150282}})     // Independence Hall
	stores = append(stores, StoreData{83, "526 Market St, Philadelphia, PA, 19106", Loc{39.949557, -75.149366}})       // Liberty Bell
	stores = append(stores, StoreData{84, "100 W Hersheypark Dr, Hershey, PA, 17033", Loc{40.288476, -76.657377}})     // Hersheypark
	stores = append(stores, StoreData{85, "1491 Mill Run Rd, Mill Run, PA, 15464", Loc{39.906374, -79.467479}})        // Fallingwater
	stores = append(stores, StoreData{86, "117 Sandusky St, Pittsburgh, PA, 15212", Loc{40.448256, -80.002636}})       // The Andy Warhol Museum
	stores = append(stores, StoreData{87, "400 Broad St, Seattle, WA, 98109", Loc{47.620506, -122.349277}})            // Space Needle
	stores = append(stores, StoreData{88, "85 Pike St, Seattle, WA, 98101", Loc{47.608316, -122.340130}})              // Pike Place Market
	stores = append(stores, StoreData{89, "3002 Mt Angeles Rd, Port Angeles, WA, 98362", Loc{48.114854, -123.430177}}) // Olympic NP (Visitor Center)
	stores = append(stores, StoreData{90, "55210 238th Ave E, Ashford, WA, 98304", Loc{46.756480, -121.814460}})       // Mount Rainier NP

	stores = append(stores, StoreData{91, "325 5th Ave N, Seattle, WA, 98109", Loc{47.621483, -122.348527}})                  // Museum of Pop Culture
	stores = append(stores, StoreData{92, "190 Marietta St NW, Atlanta, GA, 30303", Loc{33.757747, -84.394064}})              // CNN Center
	stores = append(stores, StoreData{93, "225 Baker St NW, Atlanta, GA, 30313", Loc{33.762709, -84.394364}})                 // Georgia Aquarium
	stores = append(stores, StoreData{94, "450 Auburn Ave NE, Atlanta, GA, 30312", Loc{33.755220, -84.372210}})               // Martin Luther King Jr. NHP
	stores = append(stores, StoreData{95, "1000 Robert E Lee Blvd, Stone Mountain, GA, 30083", Loc{33.811349, -84.145321}})   // Stone Mountain Park
	stores = append(stores, StoreData{96, "301 Martin Luther King Jr Blvd, Savannah, GA, 31401", Loc{32.080706, -81.098518}}) // Savannah Visitor Center
	stores = append(stores, StoreData{97, "360 Meeting St, Charleston, SC, 29403", Loc{32.791313, -79.935479}})               // Charleston Visitor Center
	stores = append(stores, StoreData{98, "40 Patriots Point Rd, Mt Pleasant, SC, 29464", Loc{32.790508, -79.907531}})        // Patriots Point
	stores = append(stores, StoreData{99, "1325 Celebrity Cir, Myrtle Beach, SC, 29577", Loc{33.714237, -78.882079}})         // Broadway at the Beach
	stores = append(stores, StoreData{100, "1401 National Park Dr, Manteo, NC, 27954", Loc{35.938022, -75.718014}})           // Fort Raleigh / Outer Banks
	stores = append(stores, StoreData{101, "1 Alamo Sq, San Francisco, CA, 94117", Loc{37.776197, -122.432530}})              // Painted Ladies (Alamo Square)
}
