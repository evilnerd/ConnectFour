
async function getGames() {
    console.log("Requesting games list")
    const res = await fetch('http://localhost:8443/game',{ cache: 'no-store' })
    // The return value is *not* serialized
    // You can return Date, Map, Set, etc.

    if (!res.ok) {
      // This will activate the closest `error.js` Error Boundary
      throw new Error('Failed to fetch data')
    }
   
    return res.json()
  }
   
 export default async function GamesList() {
    const data = await getGames()
    return <main>
        <h2 className="text-2xl font-semibold">Games</h2>
        <div id="games-list">
        { data.map((element: { key: string, created_by: string, created_at: string }) => 
              <a key={element.key} href={`join?key=${element.key}`} className="group rounded-lg border border-transparent px-5 py-4 transition-colors hover:border-gray-300 hover:bg-gray-100 hover:dark:border-neutral-700 hover:dark:bg-neutral-800/30">
                <div className="p-2 text-xl font-mono font-semibold ">{element.key}</div>
                <div>
                  <span>Created by</span>
                  <span className="p-2 font-semibold">{ element.created_by }</span>
                  <span>at</span>
                  <span className="p-2 font-semibold">{ (new Date(element.created_at)).toISOString() }</span></div>
              </a>
        )}
        </div>
    </main>
  }


  /*
  * Game create response:
  {"key":"hold alls yang","created_at":"2024-08-11T22:19:19.1940452+02:00","created_by":"Dick Appel","status":"created"}

  * Game list response:
  [{"key":"hold alls yang","created_at":"2024-08-11T22:20:27.4504738+02:00","created_by":"Dick Appel","status":"created"}]
  */