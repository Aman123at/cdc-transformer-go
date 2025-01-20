'use client'

import { useState } from 'react'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Button } from '@/components/ui/button'
import { Copy } from 'lucide-react'

type Document = {
  [key: string]: any
}

export default function CollectionAccordion({ collectionName,documents }: { collectionName: string,documents:Document[] }) {

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
  }

  return (
    <Accordion type="single" collapsible className="w-full mb-2">
      <AccordionItem value={collectionName}>
        <AccordionTrigger>{collectionName}</AccordionTrigger>
        <AccordionContent>
          <ScrollArea className="h-[200px] w-full rounded-md border p-4">
            {documents.map((doc) => (
              <div key={doc.id} className="mb-4 last:mb-0">
                <div className="flex justify-between items-start">
                  <pre className="text-sm">{JSON.stringify(doc, null, 2)}</pre>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => copyToClipboard(JSON.stringify(doc))}
                  >
                    <Copy className="h-4 w-4" />
                  </Button>
                </div>
                <hr className="my-2 border-t border-gray-200 dark:border-gray-700" />
              </div>
            ))}
          </ScrollArea>
        </AccordionContent>
      </AccordionItem>
    </Accordion>
  )
}

