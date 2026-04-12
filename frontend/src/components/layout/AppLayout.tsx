'use client'

import React from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { LayoutDashboard, ListTree, FileLock2, Activity, Settings, Search, LogOut } from 'lucide-react'
import { useAuthStore } from '@/store/useAuthStore'

const navItems = [
  { name: 'Dashboard', href: '/', icon: LayoutDashboard },
  { name: 'Topics', href: '/topics', icon: ListTree },
  { name: 'Policies', href: '/policies', icon: FileLock2 },
  { name: 'Audit Logs', href: '/audit', icon: Activity },
  { name: 'Settings', href: '/settings', icon: Settings },
]

export function AppLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const { user, logout } = useAuthStore()

  return (
    <div className="flex h-screen overflow-hidden bg-background">
      {/* Sidebar */}
      <aside className="w-64 border-r border-border bg-card hidden md:flex flex-col">
        <div className="h-14 flex items-center px-6 border-b border-border font-bold text-lg">
          Kafka<span className="text-primary ml-1">Gov</span>
        </div>
        <nav className="flex-1 overflow-y-auto py-4">
          <ul className="space-y-1 px-3">
            {navItems.map((item) => {
              const isActive = pathname === item.href || (pathname.startsWith(item.href) && item.href !== '/')
              return (
                <li key={item.name}>
                  <Link href={item.href} className={`flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors ${isActive ? 'bg-primary/10 text-primary' : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'}`}>
                    <item.icon className="h-4 w-4" />
                    {item.name}
                  </Link>
                </li>
              )
            })}
          </ul>
        </nav>
        <div className="p-4 border-t border-border">
          <div className="flex items-center gap-3 px-3 py-2">
            <div className="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center text-primary font-bold">
              {user?.username?.[0]?.toUpperCase() || 'U'}
            </div>
            <div className="flex-1 overflow-hidden text-sm">
              <p className="truncate font-medium">{user?.username || 'User'}</p>
              <p className="truncate text-xs text-muted-foreground">{user?.role || 'Guest'}</p>
            </div>
            <button onClick={logout} className="text-muted-foreground hover:text-destructive">
              <LogOut className="h-4 w-4" />
            </button>
          </div>
        </div>
      </aside>

      {/* Main Content */}
      <div className="flex-1 flex flex-col">
        {/* Topbar */}
        <header className="h-14 flex items-center justify-between px-6 border-b border-border bg-card">
          <div className="flex items-center gap-4 flex-1">
            <div className="relative w-64">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <input type="search" placeholder="Search topics or policies..." className="w-full h-9 pl-9 pr-4 text-sm rounded-md border border-input bg-background/50 focus:outline-none focus:ring-1 focus:ring-primary" />
            </div>
          </div>
          <div className="flex items-center gap-4">
            <div className="text-xs font-semibold px-2 py-1 rounded bg-secondary text-secondary-foreground border">
              Cluster: local-dev
            </div>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-y-auto p-6 bg-background">
          {children}
        </main>
      </div>
    </div>
  )
}
