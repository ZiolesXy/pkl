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
import { getBarangs } from '@/lib/api/barang'
import DeleteButton from '../DeleteButton'

async function CardTableBarang() {
  const barangs = await getBarangs()

  return (
    <Card className="w-full h-105">
      <CardHeader className="border-b">
        <CardTitle>Barang</CardTitle>
      </CardHeader>
      <CardContent className="h-85 overflow-hidden">
        <div className="h-full w-full overflow-auto ">
        <Table className="w-full table-fixed">
          <TableHeader>
            <TableRow>
              <TableHead className="w-24">ID</TableHead>
              <TableHead className="w-auto">Name</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
           {barangs.map((barang : any) => (
              <TableRow key={barang.id}>
                <TableCell className="font-medium">{barang.id}</TableCell>
                <TableCell className="truncate">{barang.name}</TableCell>
                <TableCell>
                  <DeleteButton 
                  type='barang'
                  barangId={barang.id}/>
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

export default CardTableBarang
