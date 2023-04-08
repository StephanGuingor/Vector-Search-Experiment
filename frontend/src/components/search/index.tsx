import { useState } from "react"
import { AiOutlineSearch } from "react-icons/ai"

interface SearchProps {
    placeholder: string
    onSearch: (q: string) => void
}

export function Search({ placeholder, onSearch }: SearchProps) {
    const [value, setValue] = useState('')

    function onChange(event: React.ChangeEvent<HTMLInputElement>) {
        setValue(event.target.value)
    }

  function onKeyDown(event: React.KeyboardEvent<HTMLInputElement>) {
    if (event.key === 'Enter') {
        event.preventDefault()
        onSearch(value)
        setValue('')
    }
  }

  return (
    <div className="mb-6 w-full flex rounded-full border border-slate-500 pl-2.5 pr-3 h-12 hover:shadow-lg">
        <div className="self-center flex pr-2">
            <AiOutlineSearch className="fill-slate-500"/>
        </div>
        <input onKeyDown={onKeyDown} value={value} onChange={onChange} placeholder={placeholder} autoComplete="off" type="text" id="default-input" className="border-none flex-1 outline-0 m-1 text-slate-900 text-sm block w-full dark:placeholder-slate-400 dark:text-white bg-transparent focus:ring-0 focus:ring-offset-0"/>
    </div>
  )
}
