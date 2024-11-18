import { useState } from "react"
import Search from "./Search"
import { Button } from "../ui/button";
import { useAppDispatch } from "../../redux/hooks";

interface ItemRequirements {
    id: string,
}

interface SearchableListProps<T> {
    items: T[]
    filterFunc: (searchTerm: string) => (item: T) => boolean,
    sortFunc: (item1: T, item2: T) => number,
    selectItem: (item: T) => void,
    displayFunc: (item: T) => string,
    newItemFunc: any,
    typeName: string,
}

export default function SearchableList<T extends ItemRequirements>(
    {
        items,
        filterFunc,
        sortFunc,
        selectItem,
        displayFunc,
        newItemFunc,
        typeName,
    }: SearchableListProps<T>) {

    const [searchTerm, setSearchTerm] = useState("");
    const [selectedItem, setSelectedItem] = useState<T | null>(null)
    const dispatch = useAppDispatch();

    return (
        <div className="flex flex-col items-center align-middle">
            <div className="w-1/2 ml-auto">
                <Search
                    className="w-2/3 inline-block border-b-0"
                    searchTerm={searchTerm}
                    setSearchTerm={setSearchTerm}
                />
                <Button
                    className="inline-block bg-background hover:bg-background-dark text-white w-1/3"
                    onClick={() => {
                        setSearchTerm("");
                    }}
                >
                    Clear
                </Button>
            </div>
            <ul className="overflow-x-scroll ml-auto h-2/3 w-1/2 border-black border p-2">
                {items
                    .filter(filterFunc(searchTerm))
                    .sort(sortFunc)
                    .map((item) => (
                        <li
                            key={item.id}
                            className={`${selectedItem?.id === item.id ? "bg-background text-white" : ""} cursor-pointer`}
                            onClick={() => {
                                setSelectedItem(item)
                                selectItem(item)
                            }}
                        >
                            {displayFunc(item)}
                        </li>
                    ))}
            </ul>
            <Button
                className="bg-background text-white hover:bg-background-dark ml-auto mt-2"
                onClick={() => {
                    dispatch(newItemFunc(selectedItem));
                }}
            >
                Add {typeName}
            </Button>
        </div>
    );
}