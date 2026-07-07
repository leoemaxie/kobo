<script lang="ts">
  import { enhance } from '$app/forms';

  let { data, form } = $props();
  let admins = $derived(data.admins);

</script>

<div class="space-y-10 w-full">
  <div class="flex items-center justify-between border-b border-iron pb-6">
    <div>
      <h1 class="text-3xl font-bold text-pure-white tracking-tight">Superadmin Console</h1>
      <p class="text-smoke mt-1">Manage administrators, grant access, and configure scopes.</p>
    </div>
    <div class="bg-carbon border border-electric-lime text-electric-lime text-xs font-bold px-3 py-1.5 rounded-full tracking-widest uppercase shadow-[0_0_10px_rgba(204,255,0,0.2)]">
      Superadmin
    </div>
  </div>

  <div class="space-y-6">
    {#if form?.success}
      <div class="bg-dark-olive/20 border border-electric-lime/50 text-electric-lime text-sm p-4 rounded-lg shadow-sm">
        Operation successful!
      </div>
    {/if}
    {#if form?.error}
      <div class="bg-danger/10 border border-danger/50 text-danger text-sm p-4 rounded-lg shadow-sm">
        {form.error}
      </div>
    {/if}

    <div class="bg-carbon border border-iron rounded-xl overflow-hidden shadow-sm">
      <div class="overflow-x-auto">
        <table class="min-w-full">
          <thead>
            <tr class="border-b border-iron/50 bg-void-black/50">
              <th class="px-6 py-4 text-left text-[10px] font-bold uppercase tracking-widest text-smoke">Name</th>
              <th class="px-6 py-4 text-left text-[10px] font-bold uppercase tracking-widest text-smoke">Email</th>
              <th class="px-6 py-4 text-left text-[10px] font-bold uppercase tracking-widest text-smoke">Status</th>
              <th class="px-6 py-4 text-left text-[10px] font-bold uppercase tracking-widest text-smoke">Scope</th>
              <th class="px-6 py-4 text-right text-[10px] font-bold uppercase tracking-widest text-smoke">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-iron/50">
            {#each admins as admin}
              <tr class="hover:bg-graphite/30 transition-colors">
                <td class="px-6 py-4 text-sm font-medium text-pure-white whitespace-nowrap">{admin.name}</td>
                <td class="px-6 py-4 text-sm text-smoke whitespace-nowrap">{admin.email}</td>
                <td class="px-6 py-4 text-sm whitespace-nowrap">
                  {#if admin.status === 'active'}
                    <span class="inline-flex items-center rounded-full bg-dark-olive/20 border border-electric-lime/30 px-2.5 py-0.5 text-[10px] font-bold text-electric-lime uppercase tracking-wider">Active</span>
                  {:else if admin.status === 'pending'}
                    <span class="inline-flex items-center rounded-full bg-paper/10 border border-fog/30 px-2.5 py-0.5 text-[10px] font-bold text-paper uppercase tracking-wider">Pending</span>
                  {:else}
                    <span class="inline-flex items-center rounded-full bg-danger/10 border border-danger/30 px-2.5 py-0.5 text-[10px] font-bold text-danger uppercase tracking-wider">Revoked</span>
                  {/if}
                </td>
                <td class="px-6 py-4 text-sm whitespace-nowrap">
                  <form method="POST" action="?/updateScope" use:enhance class="flex items-center gap-2">
                    <input type="hidden" name="adminId" value={admin.id} />
                    <select name="scope" class="bg-void-black border border-iron rounded text-xs text-paper px-2 py-1 outline-none focus:border-electric-lime" onchange={(e) => e.currentTarget.form?.requestSubmit()}>
                      <option value="read" selected={admin.scope === 'read'}>Read</option>
                      <option value="write" selected={admin.scope === 'write'}>Write</option>
                      <option value="full" selected={admin.scope === 'full'}>Full Access</option>
                    </select>
                  </form>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right flex justify-end gap-2">
                  {#if admin.status === 'pending' || admin.status === 'revoked'}
                    <form method="POST" action="?/updateStatus" use:enhance>
                      <input type="hidden" name="adminId" value={admin.id} />
                      <input type="hidden" name="status" value="active" />
                      <button type="submit" class="text-[11px] font-bold bg-electric-lime text-void-black hover:bg-lime-glow px-3 py-1.5 rounded transition-all shadow-sm">
                        Grant Access
                      </button>
                    </form>
                  {/if}
                  {#if admin.status === 'active'}
                    <form method="POST" action="?/updateStatus" use:enhance>
                      <input type="hidden" name="adminId" value={admin.id} />
                      <input type="hidden" name="status" value="revoked" />
                      <button type="submit" class="text-[11px] font-bold text-danger border border-danger/30 hover:bg-danger/10 px-3 py-1.5 rounded transition-colors">
                        Revoke
                      </button>
                    </form>
                  {/if}
                </td>
              </tr>
            {/each}
            {#if admins.length === 0}
              <tr>
                <td colspan="5" class="px-6 py-12 text-center text-smoke text-sm">
                  No administrators found.
                </td>
              </tr>
            {/if}
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>
