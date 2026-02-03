import { getOwnership } from "@/lib/api/ownership"
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { PostOwnershipButton } from "../PostOwnershipButton"

async function PutCardTableOwnership() {
  const ownership = await getOwnership()

  return (
    <Card className="w-full h-105">
      <CardHeader className='border-b flex justify-between'>
        <CardTitle>List of Ownership</CardTitle>
      </CardHeader>
      <CardContent className="h-85 overflow-hidden">
        <div className="h-full w-full overflow-auto">
          <Table className="w-full table-fixed">
            <TableHeader>
              <TableRow>
                <TableHead className="w-32">Barang ID</TableHead>
                <TableHead className="w-32">User ID</TableHead>
                <TableHead className="w-28">Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {ownership.map((o, idx) => (
                <TableRow key={`${o.user_id}-${o.barang_id}-${idx}`}>
                  <TableCell className="font-medium">{o.barang_id}</TableCell>
                  <TableCell>{o.user_id}</TableCell>

                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  )
}

export default PutCardTableOwnership
