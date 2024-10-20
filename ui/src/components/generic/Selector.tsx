import { Popover, PopoverTrigger } from "../ui/popover"
import { Button } from "../ui/button"
import React, { SetStateAction, useState } from "react"
import { Check, ChevronsUpDown } from "lucide-react"
import { PopoverContent } from "@radix-ui/react-popover"
import { Command, CommandEmpty, CommandGroup } from "../ui/command"
import { CommandItem, CommandList } from "cmdk"

interface SelectorProps {
    value: string
    setValue: React.Dispatch<SetStateAction<string>>
    options: { value: string, label: string }[]
}

export default function Selector({ value, setValue, options }: SelectorProps) {
    const [open, setOpen] = useState(false)
    return (
        <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
                <Button variant="outline"
                    role="combobox"
                    aria-expanded={open}
                    className="bg-secondary
                    w-40"
                >
                    {value ? options.find(option => option.value === value)?.label : "Select Option..."}
                    <ChevronsUpDown />
                </Button>
            </PopoverTrigger>
            <PopoverContent className="bg-primary w-40">
                <Command>
                    <CommandList>
                        <CommandEmpty>Not found.</CommandEmpty>
                        <CommandGroup className="bg-background">
                            {options.map(option => (
                                < CommandItem
                                    key={option.value}
                                    value={option.value}
                                    onSelect={currentValue => {
                                        setValue(currentValue)
                                        setOpen(false)
                                    }}
                                    className="cursor-pointer py-2 hover:bg-accent hover:text-black"
                                >
                                    <Check className={value === option.value ? "opacity-100 inline" : "opacity-0 inline"} /> {option.label}
                                </CommandItem>
                            )
                            )}
                        </CommandGroup>
                    </CommandList>
                </Command>
            </PopoverContent>
        </Popover >
    )
}