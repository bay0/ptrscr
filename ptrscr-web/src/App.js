import React, { useState, useEffect } from 'react';
import './assets/css/app.css'
import './assets/css/main.css'

function App() {
  let [loading, setLoading] = useState(true)
  let [image, setImage] = useState("")
  
  useEffect(() => {
    const url = new URL(window.location);
    const imgUrl = url.searchParams.get("id");
    console.log(imgUrl)
    
    if(imgUrl !== null) {
      const requestOptions = {
        method: 'GET',
        redirect: 'follow'
      };
      
      fetch("https://gist.githubusercontent.com/raw/"+imgUrl, requestOptions)
        .then(response => response.text())
        .then(result => {
          setImage(result)
          setLoading(false)
        })
    }
  }, []);

  const saveBase64AsFile = () => {
    let link = document.createElement("a");
    document.body.appendChild(link); // for Firefox
    link.setAttribute("href", image);
    link.setAttribute("download", "img.png");
    link.click();
  }

  return (
    <div className="App p-10 min-h-screen flex md:flex-row items-center justify-around bg-white dark:bg-gray-800 flex-wrap sm:flex-col font-mono">
      {loading && (
        <div className="p-0.2">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="lg:text-center">
              <div className="block">
                <h2 className="bg-white dark:bg-gray-800 text-base text-indigo-600 font-semibold tracking-wide uppercase text-5xl"><a href="https://github.com/bay0/ptrscr">PTRSCR.DEV</a></h2>
              </div>
              <p className="mt-2 text-3xl leading-8 font-extrabold tracking-tight text-white sm:text-4xl">
                Upload screenshots to GitHub gist
              </p>
              <p className="mt-4 max-w-2xl text-xl text-gray-500 lg:mx-auto">
                <a href="https://github.com/bay0/ptrscr">
                  <button
                      type="button"
                      className={`border text-4xl border-white-500 text-white rounded-md px-4 py-2 m-2 transition duration-500 ease select-none hover:text-white focus:outline-none focus:shadow-outline transform hover:-translate-y-1 hover:scale-110`}
                  > 
                    Source
                  </button>
                </a>
                <a href="https://github.com/bay0/ptrscr/releases">
                  <button
                      type="button"
                      className={`border text-4xl border-white-500 text-white rounded-md px-4 py-2 m-2 transition duration-500 ease select-none hover:text-white focus:outline-none focus:shadow-outline transform hover:-translate-y-1 hover:scale-110`}
                  > 
                    Download
                  </button>
                </a>
              </p>
            </div>
          </div>
        </div>
      )}
      {!loading && [
        <div className={`${loading ? 'animate-pulse' : ''} block p-0.2`}>
          <img src={image} width="100%" height="auto" alt="" />
        </div>,
        <div className={`${loading ? 'animate-pulse' : ''}`}>
          <button
              onClick={saveBase64AsFile}
              type="button"
              className={`border text-7xl border-white-500 text-white rounded-md px-4 py-2 m-2 transition duration-500 ease select-none hover:text-white focus:outline-none focus:shadow-outline transform hover:-translate-y-1 hover:scale-110`}
            >
              Download
          </button>
        </div>
      ]}
    </div>
  );
}

export default App;
