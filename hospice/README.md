# About cook/hospice

For some fragmentation task, you can use hospice.Go() run that.

And in main program, after main task complete, call hospice.Wait(), wait all task done

## In main flow

```
// do something

hospice.Wait()
```

## In anyway

```
hospice.Go(func(param1, param2 interface{}) {
    // do something
}, arg1, arg2)
```
