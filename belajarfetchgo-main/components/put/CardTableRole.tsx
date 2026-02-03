import { Card, CardContent, CardHeader, CardTitle } from '../ui/card'
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
import { getRoles } from '@/lib/api/roles-axios'
import PutButton from '../PutButton'

async function PutCardTableRole() {
  const roles = await getRoles()
  return (
    <Card className="w-full h-105">
      <CardHeader className='border-b flex justify-between'>
        <CardTitle>List of Role</CardTitle>
      </CardHeader>
      <CardContent className="h-85 overflow-hidden">
        <div className="h-full w-full overflow-auto">
          <Table className="w-full table-fixed">
            <TableHeader>
              <TableRow>
                <TableHead className="w-24">ID</TableHead>
                <TableHead className="w-auto">Name</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {roles.map((role: any) => (
                <TableRow key={role.id}>
                  <TableCell className="font-medium">{role.id}</TableCell>
                  <TableCell className="truncate">{role.name}</TableCell>
                  <TableCell>
                    <PutButton
                      id={role.id}
                      endpoint="/role"
                      defaultValues={{
                        name: role.name,
                      }}
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  )
}

export default PutCardTableRole
