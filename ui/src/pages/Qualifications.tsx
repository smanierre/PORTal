// import { useState } from "react"
import useLoginRequired from "../hooks/useLoginRequired"

export default function Qualifications() {
    useLoginRequired()
    // const [quals, setQuals] = useState("")
    // useEffect(() => {
    //     const getQuals = async () => {
    //         const res = await fetch(`http://localhost:8080/api/member/${memberCtx.Member?.id}/qualifications`)
    //         if (res.status !== 200) {
    //             console.log("error!!!!")
    //         }
    //         const j = await res.json()
    //         setQuals(JSON.stringify(j))
    //     }
    //     getQuals()
    // }, [])
    return <div>
        {/* {quals} */}
    </div>
}