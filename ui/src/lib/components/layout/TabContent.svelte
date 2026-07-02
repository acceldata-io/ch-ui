<script lang="ts">
  import { getGroupActiveTab } from '../../stores/tabs.svelte'
  import type { QueryTab, TableTab, DatabaseTab, DashboardTab, ModelTab } from '../../stores/tabs.svelte'
  import QueryContent from './content/QueryContent.svelte'
  import TableContent from './content/TableContent.svelte'
  import DatabaseContent from './content/DatabaseContent.svelte'
  import SavedQueries from '../../../pages/SavedQueries.svelte'
  import Settings from '../../../pages/Settings.svelte'
  import Dashboards from '../../../pages/Dashboards.svelte'
  import BrainPage from '../../../pages/Brain.svelte'
  import Admin from '../../../pages/Admin.svelte'
  import Pipelines from '../../../pages/Pipelines.svelte'
  import Telemetry from '../../../pages/Telemetry.svelte'
  import Models from '../../../pages/Models.svelte'
  import ModelContent from './content/ModelContent.svelte'
  import Home from '../../../pages/Home.svelte'

  interface Props {
    groupId: string
  }

  let { groupId }: Props = $props()

  const activeTab = $derived(getGroupActiveTab(groupId))
</script>

<div class="flex-1 min-h-0 overflow-hidden">
  {#if !activeTab}
    <div class="flex items-center justify-center h-full text-gray-400 dark:text-gray-600 text-sm">
      Open a query or select a table to get started
    </div>
  {:else if activeTab.type === 'query'}
    {#key activeTab.id}
      <QueryContent tab={activeTab as QueryTab} />
    {/key}
  {:else if activeTab.type === 'table'}
    {#key activeTab.id}
      <TableContent tab={activeTab as TableTab} />
    {/key}
  {:else if activeTab.type === 'database'}
    {#key activeTab.id}
      <DatabaseContent tab={activeTab as DatabaseTab} />
    {/key}
  {:else if activeTab.type === 'saved-queries'}
    <SavedQueries />
  {:else if activeTab.type === 'settings'}
    <Settings />
  {:else if activeTab.type === 'dashboards'}
    <Dashboards />
  {:else if activeTab.type === 'dashboard'}
    {#key activeTab.id}
      <Dashboards dashboardId={(activeTab as DashboardTab).dashboardId} />
    {/key}
  {:else if activeTab.type === 'brain'}
    <BrainPage />
  {:else if activeTab.type === 'admin'}
    <Admin />
  {:else if activeTab.type === 'pipelines'}
    <Pipelines />
  {:else if activeTab.type === 'telemetry'}
    <Telemetry />
  {:else if activeTab.type === 'model'}
    {#key activeTab.id}
      <ModelContent tab={activeTab as ModelTab} />
    {/key}
  {:else if activeTab.type === 'models'}
    <Models />
  {:else if activeTab.type === 'home'}
    <Home />
  {/if}
</div>
