import { useState } from "react";
import { MoveLeftIcon, MoveRightIcon } from "lucide-react";
import { Button } from "../ui/button";

interface ItemRequirements {
    id: string;
    name: string;
}

interface PickerProps<T> {
    pickedItems: T[]
    unpickedItems: T[]
    setPicked: (item: T) => void
    setUnpicked: (item: T) => void
    className?: string
}
export default function Picker<T extends ItemRequirements>({ pickedItems, unpickedItems, setPicked, setUnpicked, className }: PickerProps<T>) {
    const [selectedItem, setSelectedItem] = useState<T>()
    return (
        <div
            className={`${className ? className : ""} grid grid-cols-picker`}
        >
            <div className="h-full">
                <p>Selected:</p>
                <ul className="overflow-scroll border-black h-3/4 border max-h-max">
                    {pickedItems.map(item =>
                        <li key={item.id}
                            className={`${selectedItem?.id === item.id ? "bg-background text-white" : ""} cursor-pointer`}
                            onClick={() => { setSelectedItem(item) }}
                        >{item.name}</li>
                    )}
                </ul >
            </div>
            <div className="flex flex-col items-center gap-4 justify-center">
                <Button
                    type="button"
                    // disabled={selectedItem === undefined || pickedItems.indexOf(selectedItem) !== -1}
                    onClick={() => {
                        console.log(document.activeElement)
                        if (selectedItem) {
                            setPicked(selectedItem as T)
                        }
                    }}
                >
                    <MoveLeftIcon />
                </Button>
                <Button
                    type="button"
                    // disabled={selectedItem === undefined || unpickedItems.indexOf(selectedItem) !== -1}
                    onClick={() => {
                        if (selectedItem) {
                            setUnpicked(selectedItem as T)
                        }
                    }}
                >
                    <MoveRightIcon />
                </Button>
            </div>
            <div className="h-full">
                <p>Available:</p>
                <ul className="overflow-scroll border-black border h-3/4 max-h-max">
                    {unpickedItems.map(item =>
                        <li key={item.id}
                            className={`${selectedItem?.id === item.id ? "bg-background text-white" : ""} cursor-pointer`}
                            onClick={() => {
                                setSelectedItem(item)
                            }}
                        >{item.name}</li>
                    )}
                </ul>
            </div>
        </div>
    )
}