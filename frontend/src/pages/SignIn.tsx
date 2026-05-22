import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { Link } from "react-router-dom"
import { z } from "zod"

import { AuthCard } from "@/components/AuthCard"
import { CircleIcon } from "@/components/CircleIcon"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

const schema = z.object({
  email: z.string().email("hmm, that doesn't look like an email"),
  password: z.string().min(1, "password is required"),
  remember: z.boolean().optional(),
})

type FormValues = z.infer<typeof schema>

const labelClass =
  "font-mono text-[11px] uppercase tracking-wider text-ink-mute"

export function SignIn() {
  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { email: "", password: "", remember: false },
  })

  function onSubmit(values: FormValues) {
    // submit wiring lands in Phase 10
    console.log("sign in", values)
  }

  return (
    <AuthCard>
      <div className="text-center">
        <CircleIcon className="mx-auto mb-3.5" size={56}>
          ◎
        </CircleIcon>
        <h1 className="text-[2rem] leading-tight">welcome back</h1>
        <p className="mt-1.5 text-[15px] text-ink-mute font-hand">
          sign in to keep going
        </p>
      </div>

      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmit)}
          className="flex flex-col gap-4"
        >
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem className="gap-1.5">
                <FormLabel className={labelClass}>email</FormLabel>
                <FormControl>
                  <Input
                    placeholder="you@somewhere.com"
                    autoComplete="email"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem className="gap-1.5">
                <FormLabel className={labelClass}>password</FormLabel>
                <FormControl>
                  <Input
                    type="password"
                    placeholder="••••••••"
                    autoComplete="current-password"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="flex items-center justify-between text-[13px] text-ink-mute">
            <div className="flex items-center gap-2">
              <Checkbox
                id="remember"
                checked={form.watch("remember")}
                onCheckedChange={(v) =>
                  form.setValue("remember", v === true, { shouldDirty: true })
                }
              />
              <Label htmlFor="remember" className="font-hand cursor-pointer">
                remember me
              </Label>
            </div>
            <button
              type="button"
              className="link-dashed text-ink-mute font-hand"
              onClick={() => console.log("forgot password — not wired yet")}
            >
              forgot?
            </button>
          </div>

          <Button type="submit">sign in →</Button>
        </form>
      </Form>

      <p className="text-center text-sm text-ink-mute font-hand">
        new here?{" "}
        <Link to="/sign-up" className="link-dashed text-ink">
          make an account
        </Link>
      </p>
    </AuthCard>
  )
}
