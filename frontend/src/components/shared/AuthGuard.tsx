'use client'

import { useAuthStore } from '@/store/useAuthStore'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'

export function AuthGuard({ children }: { children: React.ReactNode }) {
  const token = useAuthStore((state) => state.token)
  const router = useRouter()
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
    if (!token) {
      router.push('/login')
    }
  }, [token, router])

  if (!mounted) return null
  if (!token) return null

  return <>{children}</>
}
