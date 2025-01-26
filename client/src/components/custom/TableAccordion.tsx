'use client'

import { useState } from 'react'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Edit, Trash, Plus } from 'lucide-react'
import EditColumnModal from './EditColumnModal'
import InsertRowModal from './InsertRowModal'
import { IColumn, IRow } from '@/interfaces/commonInterface'
import { usePG } from '@/contexts/pgcontext'

interface ITableAccordionProps {
  tableName: string, columns: IColumn[],rows:IRow[]
}

export default function TableAccordion({ tableName, columns,rows }: ITableAccordionProps) {
  const [editingColumn, setEditingColumn] = useState<IColumn | null>(null)
  const [isInsertModalOpen, setIsInsertModalOpen] = useState<boolean>(false)

  const {handleInsertRow,handleDeleteRow} = usePG()

  const handleEditColumn = (column: IColumn) => {
    setEditingColumn(column)
  }

  const handleUpdateColumn = (updatedColumn: IColumn) => {
    // Update column logic here
  }

  const deleteRow = (rowid: number) => {
    // Delete column logic here
    handleDeleteRow(tableName,rowid)
  }

  const insertRow = (newRow: IRow) => {
    // call insert row api 
    handleInsertRow(tableName,newRow)
    setIsInsertModalOpen(false)
  }

  return (
    <Accordion type="single" collapsible className="w-full mb-2">
      <AccordionItem value={tableName}>
        <AccordionTrigger>{tableName}</AccordionTrigger>
        <AccordionContent>
          <Button onClick={() => setIsInsertModalOpen(true)} className="mb-2">
            <Plus className="mr-2 h-4 w-4" /> Insert a row
          </Button>
          <div className="h-[200px] overflow-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  {columns.map((column:IColumn) => (
                    <TableHead key={column.name}>{column.name}</TableHead>
                  ))}
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {rows.map((row:IRow, rowIndex:number) => (
                  <TableRow key={rowIndex}>
                    {columns.map((column:IColumn) => (
                      <TableCell key={column.name}>{row[column.name]}</TableCell>
                    ))}
                    <TableCell>
                      {/* <Button variant="ghost" size="icon" onClick={() => handleEditColumn(columns[rowIndex])}>
                        <Edit className="h-4 w-4" />
                      </Button> */}
                      <Button variant="ghost" size="icon" onClick={() => deleteRow(row.id)}>
                        <Trash className="h-4 w-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </AccordionContent>
      </AccordionItem>
      {editingColumn && (
        <EditColumnModal
          isOpen={!!editingColumn}
          onClose={() => setEditingColumn(null)}
          column={editingColumn}
          onUpdateColumn={handleUpdateColumn}
        />
      )}
      <InsertRowModal
        isOpen={isInsertModalOpen}
        onClose={() => setIsInsertModalOpen(false)}
        columns={columns}
        onInsertRow={insertRow}
      />
    </Accordion>
  )
}

