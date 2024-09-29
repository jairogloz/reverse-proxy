import { useEffect, useState } from 'react'


export function useFetch(url) {
    const [data, setData] = useState(0);
    const [loading, seLoading] = useState(true);
    const [error, setError] = useState(null);
    const [controller, setController] = useState(null)
  

    useEffect(() => {
        const abortController = new AbortController();
        setController(abortController)
        seLoading(true)
        fetch(url, {signal: abortController.signal})
            .then((response) => response.json())
            .then((data)=> setData(data))
            .catch((err) => {
                if (err.name === "AbortError"){
                    console.log("Request cancelled");
                }else {
                    setError(err)
                }
            })
            .finally(()=> seLoading(false))
        return () => abortController.abort();
    }, [] );

    const handleCancelRequest = () => {
        if (controller) {
            controller.abort();
            setError("Request cancelled")
        }
    }
    

    return { data, loading, error, handleCancelRequest };
}