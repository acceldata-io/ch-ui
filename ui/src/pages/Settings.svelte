<script lang="ts">
  import { getSession } from '../lib/stores/session.svelte'
  import {
    Settings as SettingsIcon,
    FileText,
    ExternalLink,
    Building2,
  } from 'lucide-svelte'
  import logo from '../assets/logo.png'

  const session = $derived(getSession())
</script>

<div class="h-full overflow-auto">
  <div class="max-w-7xl mx-auto p-6 space-y-4">
    <section class="ds-panel overflow-hidden">
      <div class="px-5 py-4 border-b border-orange-300/25 dark:border-orange-700/35 bg-gradient-to-r from-orange-100/60 via-transparent to-transparent dark:from-orange-500/10 dark:via-transparent">
        <div class="flex items-center gap-4 min-w-0">
          <div class="h-16 w-16 rounded-2xl border border-orange-300/40 dark:border-orange-700/40 bg-gray-100 dark:bg-gray-900 grid place-items-center overflow-hidden shrink-0">
            <img src={logo} alt="CH-UI logo" class="h-12 w-12 object-contain" />
          </div>
          <div class="min-w-0">
            <div class="flex items-center gap-2 text-gray-900 dark:text-gray-100">
              <SettingsIcon size={18} class="text-ch-orange" />
              <h1 class="text-2xl font-semibold leading-tight">Settings</h1>
            </div>
            <p class="mt-1 text-sm text-gray-600 dark:text-gray-400">Instance and connection details</p>
          </div>
        </div>
      </div>

      <div class="p-5">
        <div class="flex items-center gap-2 mb-3">
          <Building2 size={16} class="text-ch-orange" />
          <h2 class="text-sm font-semibold text-gray-800 dark:text-gray-200 uppercase tracking-wide">Brand & Instance</h2>
        </div>

        <div class="flex items-center gap-3">
          <img src={logo} alt="CH-UI mark" class="h-12 w-12 rounded-xl border border-orange-300/35 dark:border-orange-700/35 bg-gray-100 dark:bg-gray-900 p-1.5" />
          <div>
            <p class="text-sm font-semibold text-gray-900 dark:text-gray-100">CH-UI</p>
            <p class="text-xs text-gray-500 dark:text-gray-400">ClickHouse operations and analytics control surface</p>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-2 mt-3">
          <div class="ds-panel-muted p-2.5">
            <p class="text-[11px] uppercase tracking-wide text-gray-500 dark:text-gray-400">Connected As</p>
            <p class="text-xs font-medium text-gray-800 dark:text-gray-200 mt-1">{session?.user || '—'}</p>
          </div>
          <div class="ds-panel-muted p-2.5">
            <p class="text-[11px] uppercase tracking-wide text-gray-500 dark:text-gray-400">Role</p>
            <p class="text-xs font-medium text-gray-800 dark:text-gray-200 mt-1 uppercase">{session?.role || '—'}</p>
          </div>
          <div class="ds-panel-muted p-2.5">
            <p class="text-[11px] uppercase tracking-wide text-gray-500 dark:text-gray-400">Connection</p>
            <p class="text-xs font-medium text-gray-800 dark:text-gray-200 mt-1">{session?.connectionName || '—'}</p>
          </div>
          <div class="ds-panel-muted p-2.5">
            <p class="text-[11px] uppercase tracking-wide text-gray-500 dark:text-gray-400">Version</p>
            <p class="text-xs font-medium text-gray-800 dark:text-gray-200 mt-1">{session?.version || '—'}</p>
          </div>
        </div>

        <div class="mt-4 pt-3 border-t border-gray-200 dark:border-gray-800 flex items-center gap-2 flex-wrap">
          <a
            href="https://www.apache.org/licenses/LICENSE-2.0"
            target="_blank"
            rel="noreferrer"
            class="inline-flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 hover:text-ch-orange"
          >
            <FileText size={11} /> Apache-2.0 <ExternalLink size={10} />
          </a>
          <a
            href="https://github.com/caioricciuti/ch-ui/blob/main/docs/legal/terms-of-service.md"
            target="_blank"
            rel="noreferrer"
            class="inline-flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 hover:text-ch-orange"
          >
            Terms <ExternalLink size={10} />
          </a>
          <a
            href="https://github.com/caioricciuti/ch-ui/blob/main/docs/legal/privacy-policy.md"
            target="_blank"
            rel="noreferrer"
            class="inline-flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 hover:text-ch-orange"
          >
            Privacy <ExternalLink size={10} />
          </a>
        </div>
      </div>
    </section>
  </div>
</div>
