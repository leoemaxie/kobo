<script lang="ts">
  import { ShieldCheck, UploadCloud, FileText, CheckCircle2, AlertTriangle } from '@lucide/svelte';
  import { useConsoleState } from '$lib/state/console.svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import CardSection from '$lib/components/ui/CardSection.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Button from '$lib/components/ui/Button.svelte';

  const state = useConsoleState();
</script>

<div class="flex flex-col gap-7">
  <div class="mb-6 sm:mb-8">
    <div
      class="p-3 sm:p-4 mt-4 mb-3 rounded-lg bg-[var(--accent)]/10 border border-[var(--accent)]/20 flex items-start sm:items-center gap-3 text-main"
    >
      <AlertTriangle class="w-4 h-4 sm:w-5 sm:h-5 text-[var(--accent)] shrink-0 mt-0.5 sm:mt-0" />
      <p class="text-[12px] sm:text-sm font-medium leading-snug">
        <strong>Notice:</strong>
        <span class="sm:hidden">Submissions are temporarily disabled.</span>
        <span class="hidden sm:inline"
          >KYC verification is currently in development. Submissions are temporarily disabled.</span
        >
      </p>
    </div>
    <p class="text-[13px] sm:text-[14px] text-muted leading-relaxed">
      <span class="sm:hidden">Complete compliance verification to unlock live transactions.</span>
      <span class="hidden sm:inline"
        >Complete your compliance verification to unlock production access and start processing live
        transactions.</span
      >
    </p>
  </div>

  <form
    class="opacity-50 pointer-events-none select-none grayscale-[0.5]"
    method="POST"
    action="?/submitKyc"
    use:enhance={() => {
      return async ({ result, update }) => {
        if (result.type === 'failure' || result.type === 'error') {
          toast.error('Failed to submit KYC details.');
        } else {
          toast.success('KYC submitted successfully. Verification pending.');
        }
        await update();
      };
    }}
  >
    <div class="grid gap-6">
      <CardSection
        title="Business Information"
        subtitle="Details about your registered business entity."
      >
        <div class="grid gap-4 md:grid-cols-2">
          <Input
            id="businessName"
            label="Legal Business Name"
            type="text"
            name="businessName"
            required
            variant="settings"
          />
          <Input
            id="regNumber"
            label="Registration Number (RC)"
            type="text"
            name="regNumber"
            required
            variant="settings"
          />
          <div class="md:col-span-2">
            <Input
              id="address"
              label="Business Address"
              type="text"
              name="address"
              required
              variant="settings"
            />
          </div>
        </div>
      </CardSection>

      <CardSection
        title="Director Details"
        subtitle="Primary point of contact for compliance matters."
      >
        <div class="grid gap-4 md:grid-cols-2">
          <Input
            id="directorName"
            label="Full Name"
            type="text"
            name="directorName"
            required
            variant="settings"
          />
          <Input
            id="bvn"
            label="Bank Verification Number (BVN)"
            type="text"
            name="bvn"
            required
            variant="settings"
          />
        </div>
      </CardSection>

      <CardSection
        title="Document Upload"
        subtitle="Upload supporting documents. PDF or Image format, max 5MB."
      >
        <div class="grid gap-4 md:grid-cols-2">
          <label
            class="border-2 border-dashed border-[var(--border-color)] rounded-xl p-6 flex flex-col items-center justify-center text-center hover:bg-[var(--bg-element)] transition-colors cursor-pointer group"
          >
            <div
              class="p-3 bg-[var(--bg-active)] rounded-full mb-3 group-hover:bg-[var(--accent)]/10 transition-colors"
            >
              <FileText
                class="w-6 h-6 text-muted group-hover:text-[var(--accent)] transition-colors"
              />
            </div>
            <p class="text-sm font-medium text-main mb-1">Certificate of Incorporation</p>
            <p class="text-xs text-subtle">Click to browse or drag and drop</p>
            <input type="file" name="certOfInc" class="hidden" accept=".pdf,.png,.jpg,.jpeg" />
          </label>

          <label
            class="border-2 border-dashed border-[var(--border-color)] rounded-xl p-6 flex flex-col items-center justify-center text-center hover:bg-[var(--bg-element)] transition-colors cursor-pointer group"
          >
            <div
              class="p-3 bg-[var(--bg-active)] rounded-full mb-3 group-hover:bg-[var(--accent)]/10 transition-colors"
            >
              <UploadCloud
                class="w-6 h-6 text-muted group-hover:text-[var(--accent)] transition-colors"
              />
            </div>
            <p class="text-sm font-medium text-main mb-1">Valid ID (Director)</p>
            <p class="text-xs text-subtle">Click to browse or drag and drop</p>
            <input type="file" name="validId" class="hidden" accept=".pdf,.png,.jpg,.jpeg" />
          </label>
        </div>
      </CardSection>

      <div class="flex flex-col-reverse sm:flex-row justify-end gap-3 mt-4 pt-6">
        <Button variant="ghost" type="button" class="w-full sm:w-auto">Cancel</Button>
        <Button variant="primary" type="submit" class="w-full sm:w-auto">
          <CheckCircle2 size={16} class="mr-2" /> Submit for Verification
        </Button>
      </div>
    </div>
  </form>
</div>
