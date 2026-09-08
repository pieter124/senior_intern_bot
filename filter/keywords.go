package filter

// RoleAxis answers "is this an internship or student-level position?"
// The out-of-scope list holds seniority markers only:
// nothing about function or team.
var RoleAxis = axis{
	inScope: []string{
		"intern", "interns", "internship", "internships",
		"placement", "industrial placement", "placement year",
		"summer analyst", "summer associate", "summer",
		"co op", "coop",
		"apprentice", "apprenticeship", "apprenticeships",
		"trainee", "traineeship",
		"student", "students", "undergraduate", "undergrad",
		"early career", "early careers", "campus",
	},
	outOfScope: []string{
		"senior", "staff", "principal", "director",
		"vp", "vice president", "chief", "executive",
		"head of", "manager", "architect",
	},
}

// RegionAxis answers "is this posting in EMEA?"
// Matched against Location.Name only, which Greenhouse returns as free text
// ("London, UK", "Remote", "Multiple Locations"). A location naming both an
// EMEA and a non-EMEA site contradicts itself and evaluates to unknown,
// which routes to REVIEW rather than dropping a role open in both.
var RegionAxis = axis{
	inScope: []string{
		"emea", "europe", "european",

		"uk", "united kingdom", "britain", "england", "scotland", "wales",
		"ireland", "france", "germany", "deutschland", "netherlands",
		"nederland", "belgium", "spain", "españa", "portugal", "italy",
		"italia", "switzerland", "schweiz", "suisse", "austria",
		"österreich", "poland", "polska", "czechia", "czech republic",
		"sweden", "sverige", "norway", "norge", "denmark", "danmark",
		"finland", "suomi", "estonia", "latvia", "lithuania", "romania",
		"hungary", "greece", "turkey", "türkiye",

		"london", "manchester", "birmingham", "leeds", "bristol",
		"cambridge", "oxford", "reading", "edinburgh", "glasgow",
		"belfast", "dublin", "cork", "galway",
		"paris", "lyon", "toulouse", "grenoble", "sophia antipolis",
		"berlin", "munich", "münchen", "hamburg", "frankfurt",
		"cologne", "köln", "düsseldorf", "dusseldorf", "stuttgart",
		"nuremberg", "nürnberg", "dresden", "leipzig", "karlsruhe",
		"amsterdam", "utrecht", "eindhoven", "rotterdam", "the hague",
		"brussels", "bruxelles", "antwerp", "antwerpen", "ghent", "leuven",
		"madrid", "barcelona", "valencia", "sevilla", "málaga", "malaga",
		"lisbon", "lisboa", "porto", "braga",
		"milan", "milano", "rome", "roma", "turin", "torino", "bologna",
		"zurich", "zürich", "geneva", "genève", "lausanne", "basel",
		"vienna", "wien", "graz", "linz",
		"prague", "praha", "brno", "bratislava",
		"warsaw", "warszawa", "krakow", "kraków", "wroclaw", "wrocław",
		"poznan", "poznań", "gdansk", "gdańsk", "lodz", "łódź",
	},
	outOfScope: []string{
		"usa", "united states", "canada", "mexico", "brazil", "brasil",
		"argentina", "chile", "colombia", "peru",
		"india", "china", "japan", "korea", "singapore", "malaysia",
		"australia", "new zealand", "indonesia", "thailand", "vietnam",
		"philippines", "taiwan",

		"new york", "brooklyn", "manhattan", "san francisco", "palo alto",
		"mountain view", "menlo park", "san jose", "san mateo",
		"sunnyvale", "santa clara", "cupertino", "oakland", "berkeley",
		"seattle", "bellevue", "redmond", "boston", "somerville",
		"austin", "dallas", "houston", "san antonio", "chicago",
		"denver", "boulder", "atlanta", "miami", "orlando", "tampa",
		"los angeles", "santa monica", "san diego", "irvine", "portland",
		"philadelphia", "pittsburgh", "phoenix", "tempe", "detroit",
		"ann arbor", "minneapolis", "madison", "columbus", "cleveland",
		"nashville", "charlotte", "raleigh", "durham", "salt lake city",
		"las vegas", "kansas city", "st louis", "baltimore",
		"toronto", "vancouver", "montreal", "montréal", "ottawa",
		"calgary", "waterloo",
		"sao paulo", "são paulo", "rio de janeiro", "mexico city",
		"guadalajara", "buenos aires", "bogota", "bogotá", "santiago",
		"bangalore", "bengaluru", "hyderabad", "mumbai", "chennai",
		"pune", "kolkata", "gurgaon", "gurugram", "noida", "new delhi",
		"ahmedabad", "jaipur",
		"tokyo", "osaka", "kyoto", "yokohama", "seoul", "busan",
		"shanghai", "beijing", "shenzhen", "guangzhou", "hangzhou",
		"hong kong", "taipei", "kuala lumpur", "jakarta", "bangkok",
		"manila", "ho chi minh", "hanoi",
		"sydney", "melbourne", "brisbane", "perth", "adelaide",
		"canberra", "auckland", "wellington",
	},
}

// DisciplineAxis answers "is this a software engineering or CS role?"
// This axis is used as a veto rather than a requirement (see Classify):
// an out-of-scope discipline rejects, but the absence of any discipline
// signal does not block ACCEPT.
var DisciplineAxis = axis{
	inScope: []string{
		"software", "software engineer",
		"engineer", "engineering", "developer", "development",

		"backend", "back end", "frontend", "front end", "fullstack",
		"full stack", "web", "mobile", "ios", "android",

		"systems", "platform", "infrastructure", "devops", "sre",
		"site reliability", "cloud", "distributed", "networking",
		"embedded", "firmware", "compiler", "kernel", "database",

		"data engineer", "data engineering", "data science",
		"data scientist", "machine learning", "deep learning",
		"computer vision", "robotics",

		"security engineering", "security engineer", "cybersecurity",
		"cyber security", "cryptography", "application security",

		"technology", "technical", "quantitative", "quant",
		"test engineer", "automation", "product", "product engineer",
	},
	outOfScope: []string{
		"marketing", "content marketing", "brand", "advertising",
		"public relations", "communications", "copywriter", "copywriting",
		"social media", "journalism", "editorial",

		"sales", "account executive", "account manager",
		"business development", "customer success", "customer support",
		"customer service",

		"recruiting", "recruitment", "recruiter", "talent acquisition",
		"human resources", "people operations", "payroll",

		"legal", "counsel", "paralegal", "compliance officer",

		"accounting", "audit", "auditor", "tax", "bookkeeping",
		"treasury", "procurement", "supply chain", "logistics",
		"warehouse", "facilities", "office manager", "receptionist",

		"graphic design", "interior design", "fashion",
		"nursing", "clinical", "physician", "veterinary",
		"teaching", "tutor", "barista", "retail",
	},
}
