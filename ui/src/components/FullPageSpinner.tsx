import { LoadingSpinner } from "./ui/spinner"

export default function FullPageSpinner() {
    return <div className="w-full h-full flex items-center justify-center">
        <LoadingSpinner className="w-72 h-72" />
    </div>
}