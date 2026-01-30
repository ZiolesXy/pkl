import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import getUsers from '@/lib/api/users-native'
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
async function CardTable() {
  const users = await getUsers()

  return (
    <div className="flex justify-center">
      <Card className="w-full max-w-sm">
        <CardHeader className='border-b'>
          <CardTitle>List of Users</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className='font-bold'>ID</TableHead>
                <TableHead className='font-bold'>Name</TableHead>
                <TableHead className='font-bold'>Role</TableHead>
                <TableHead className='font-bold'>Barang</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {users.map((u: any, idx: number) => (
                <TableRow key={u.id ?? idx}>
                  <TableCell>{u.id}</TableCell>
                  <TableCell>{u.name}</TableCell>
                  <TableCell>{u.role?.name}</TableCell>
                  <TableCell>{u.barangs?.map((b: any) => b.name).join(".,")}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>  
        </CardContent>
      </Card>
    </div>
  )
}

export default CardTable
