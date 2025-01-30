import * as React from "react"
import {useEffect} from "react"
import {format} from "date-fns"
import {CalendarIcon} from "lucide-react"
import {cn} from "@/lib/utils"
import {Button} from "@/components/ui/button"
import {Calendar} from "@/components/ui/calendar"
import {Popover, PopoverContent, PopoverTrigger,} from "@/components/ui/popover"

interface IProps {
    name?: string
    onSelect?: (date: Date) => void
    selected?: Date
}

export function DatePicker(props: IProps) {
    const [date, setDate] = React.useState<Date | undefined>(props.selected || undefined)
    const ref = React.useRef(null)

    useEffect(() => {
        if (ref.current) {
            const input = ref.current as HTMLInputElement
            input.form?.addEventListener("reset", () => {
                setDate(undefined)
            })
        }
    }, [ref])

    useEffect(() => {
        if (date) {
            props.onSelect && props.onSelect(date)
        }
    }, [date])

    return (
        <>
            <input type="hidden" ref={ref} name={props?.name} value={date ? date.toISOString() : ""}/>
            <Popover>
                <PopoverTrigger asChild>
                    <Button
                        variant={"outline"}
                        className={cn(
                            "min-w-[240px] justify-start text-left font-normal",
                            !date && "text-muted-foreground"
                        )}
                    >
                        <CalendarIcon/>
                        {date ? format(date, "PPP") : <span>Pick a date</span>}
                    </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                    <Calendar
                        mode="single"
                        selected={date}
                        onSelect={setDate}
                        initialFocus
                    />
                </PopoverContent>
            </Popover>
        </>
    )
}
