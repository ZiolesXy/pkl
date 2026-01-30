import { Bell, CircleFadingArrowUp, HatGlasses, Inbox, Link, Logs, ReceiptText, Settings, Trash, Upload, User } from "lucide-react"
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList, CommandSeparator } from "./ui/command"
import UserItem from "./UserItem"
function Sidebar() {
    const menuList = [
        {
            group : 'General',
            items : [
            {
                link : './',
                icon : <User />,
                text : 'Profile'
            },
            {
                link : '/',
                icon :<Trash />,
                text : 'Delete'
            },
            {
                link : '/',
                icon : <Upload />,
                text : 'Post'
            },
            {
                link : '/',
                icon : <CircleFadingArrowUp />,
                text : 'Put'
            }
        ]
        },
        {
            group : 'Setting',
            items : [
            {
                link : './',
                icon : <Settings />,
                text : 'General Setting'
            },
            {
                link : '/',
                icon : <HatGlasses />,
                text : 'Privacy'
            },
            {
                link : '/',
                icon : <Logs />,
                text : 'logs'
            }
        ]
        }
    ]
    return (
        
        <div className="sticky top-0 flex flex-col w-75 min-w-75 p-4 h-screen">
            <UserItem />
            <div className="grow">
                <Command style={{overflow : "visible"}}>
                    <CommandList style={{overflow : "visible"}}>
                        {menuList.map((menu : any, key : number) =>(
                        <CommandGroup key={key} heading={menu.group}>
                            {menu.items.map((option: any, optionKey : number) =>
                            <CommandItem key={optionKey} className="flex gap-2 cursor-pointer">
                                {option.icon}
                                {option.text}
                            </CommandItem>)}
                        </CommandGroup>))}
                    </CommandList>
                </Command>
            </div>
            <div>Settings</div>
        </div>
    )
}

export default Sidebar
