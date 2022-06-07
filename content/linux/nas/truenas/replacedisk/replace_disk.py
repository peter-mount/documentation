#!/usr/bin/env python3
import argparse
import sys

from middlewared.client import Client


def main():

    parser = argparse.ArgumentParser()
    parser.add_argument('pool')
    parser.add_argument('label')
    parser.add_argument('disk')
    parser.add_argument('passphrase', nargs='?')

    args = parser.parse_args()

    with Client() as c:
        assert '13.0-RELEASE' in c.call('system.version')
        pool = c.call('pool.query', [['name', '=', args.pool]])
        if not pool:
            print('Pool not found.')
            sys.exit(1)

        disk = list(filter(lambda x: x['devname'] == args.disk, c.call('disk.get_unused')))

        if not disk:
            print('Unused disk not found.')
            sys.exit(1)

        arg = {'label': args.label, 'disk': disk[0]['identifier']}

        if pool[0]['encrypt'] == 2:
            if not args.passphrase:
                print('Passphrse required.')
                sys.exit(1)
            args['passphrase'] = args.passphrase

        c.call('pool.replace', pool[0]['id'], arg, job=True)
        print('Replace initiated.')


if __name__ == '__main__':
    main()
