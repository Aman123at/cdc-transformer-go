'use client'
import { Button } from '@/components/ui/button'
import { ClipboardList } from 'lucide-react'
import Image from 'next/image'
import logo from '../../../assets/CDC_LOGO.png'

export default function Header({ onCheckLogs }: { onCheckLogs: () => void }) {
  return (
    <header className="flex justify-between items-center p-4 bg-secondary">
      <div className="flex items-center">
        <div className="w-10 h-10 bg-primary rounded-full mr-2">
          <Image src={logo} alt='logo' />
        </div>
        <span className="text-xl font-bold">CDC Service</span>
      </div>
      <Button onClick={onCheckLogs} variant="outline">
        <ClipboardList className="mr-2 h-4 w-4" />
        Check Transform Logs
      </Button>
    </header>
  )
}

